package lifullhomes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var allowedHosts = map[string]bool{"homes.co.jp": true, "www.homes.co.jp": true}

type Client struct {
	httpClient  *http.Client
	userAgent   string
	minInterval time.Duration
	mu          sync.Mutex
	lastRequest time.Time
	robots      map[string]robotsRules
}

type robotsRules struct{ disallow []string }

func NewClient(httpClient *http.Client, userAgent string, minInterval time.Duration) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 20 * time.Second}
	}
	if userAgent == "" {
		userAgent = "room-finder/0.1 (+local crawler)"
	}
	return &Client{httpClient: httpClient, userAgent: userAgent, minInterval: minInterval, robots: make(map[string]robotsRules)}
}

func (c *Client) Crawl(ctx context.Context, req SearchRequest) ([]RentalOffer, error) {
	refs, err := c.Search(ctx, req)
	if err != nil {
		return nil, err
	}
	offers := make([]RentalOffer, 0, len(refs))
	for _, ref := range refs {
		offer, err := c.FetchOffer(ctx, ref)
		if err != nil {
			return offers, err
		}
		offers = append(offers, offer)
	}
	return offers, nil
}

func (c *Client) Search(ctx context.Context, req SearchRequest) ([]OfferReference, error) {
	doc, err := c.getDocument(ctx, req.URL)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	var refs []OfferReference
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key != "href" {
					continue
				}
				absolute, ok := resolveURL(req.URL, attr.Val)
				if !ok || !isOfferURL(absolute) {
					continue
				}
				id := sourceOfferID(absolute)
				if id == "" {
					continue
				}
				key := id + "\x00" + absolute
				if _, exists := seen[key]; exists {
					continue
				}
				seen[key] = struct{}{}
				refs = append(refs, OfferReference{SourceOfferID: id, URL: absolute})
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(doc)
	return refs, nil
}

func (c *Client) FetchOffer(ctx context.Context, ref OfferReference) (RentalOffer, error) {
	doc, err := c.getDocument(ctx, ref.URL)
	if err != nil {
		return RentalOffer{}, err
	}
	offer := parseOffer(doc, ref)
	offer.CapturedAt = time.Now().UTC()
	return offer, nil
}

func (c *Client) getDocument(ctx context.Context, rawURL string) (*html.Node, error) {
	u, err := validateURL(rawURL)
	if err != nil {
		return nil, &CrawlError{URL: rawURL, Operation: "validate", Err: err}
	}
	if err := c.checkRobots(ctx, u); err != nil {
		return nil, &CrawlError{URL: rawURL, Operation: "robots", Err: err}
	}
	if err := c.wait(ctx); err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, &CrawlError{URL: rawURL, Operation: "request", Err: err}
	}
	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &CrawlError{URL: rawURL, Operation: "request", Retryable: true, Err: err}
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &CrawlError{URL: rawURL, Operation: "response", Retryable: resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500, Err: fmt.Errorf("unexpected HTTP status %d", resp.StatusCode)}
	}
	doc, err := html.Parse(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return nil, &CrawlError{URL: rawURL, Operation: "parse", Err: err}
	}
	return doc, nil
}

func (c *Client) wait(ctx context.Context) error {
	c.mu.Lock()
	delay := time.Until(c.lastRequest.Add(c.minInterval))
	if delay < 0 {
		delay = 0
	}
	c.lastRequest = time.Now().Add(delay)
	c.mu.Unlock()
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (c *Client) checkRobots(ctx context.Context, u *url.URL) error {
	c.mu.Lock()
	rules, cached := c.robots[u.Host]
	c.mu.Unlock()
	if !cached {
		if err := c.wait(ctx); err != nil {
			return err
		}
		robotsURL := "https://" + u.Host + "/robots.txt"
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsURL, nil)
		if err != nil {
			return err
		}
		req.Header.Set("User-Agent", c.userAgent)
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return err
		}
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		_ = resp.Body.Close()
		if readErr != nil {
			return readErr
		}
		if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("robots.txt returned HTTP status %d", resp.StatusCode)
		}
		rules = parseRobots(body)
		c.mu.Lock()
		c.robots[u.Host] = rules
		c.mu.Unlock()
	}
	for _, prefix := range rules.disallow {
		if prefix != "" && strings.HasPrefix(u.EscapedPath(), prefix) {
			return errors.New("URL is disallowed by robots.txt")
		}
	}
	return nil
}

func parseRobots(body []byte) robotsRules {
	var active bool
	var rules robotsRules
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.ToLower(strings.TrimSpace(parts[0])), strings.TrimSpace(parts[1])
		switch key {
		case "user-agent":
			active = value == "*"
		case "disallow":
			if active {
				rules.disallow = append(rules.disallow, value)
			}
		}
	}
	return rules
}

func validateURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme != "https" || !allowedHosts[strings.ToLower(u.Hostname())] {
		return nil, errors.New("URL must be an HTTPS LIFULL HOME'S URL")
	}
	return u, nil
}

func resolveURL(base, href string) (string, bool) {
	baseURL, err := validateURL(base)
	if err != nil {
		return "", false
	}
	hrefURL, err := url.Parse(strings.TrimSpace(href))
	if err != nil {
		return "", false
	}
	resolved := baseURL.ResolveReference(hrefURL)
	if _, err := validateURL(resolved.String()); err != nil {
		return "", false
	}
	return resolved.String(), true
}

func isOfferURL(rawURL string) bool {
	u, err := validateURL(rawURL)
	return err == nil && strings.Contains(u.Path, "/chintai/")
}

func sourceOfferID(rawURL string) string {
	u, err := validateURL(rawURL)
	if err != nil {
		return ""
	}
	base := path.Base(strings.TrimSuffix(u.Path, "/"))
	if base == "." || base == "/" || base == "" {
		return ""
	}
	return base
}

func parseOffer(doc *html.Node, ref OfferReference) RentalOffer {
	offer := RentalOffer{Source: Source, SourceOfferID: ref.SourceOfferID, SourceURL: ref.URL, Details: map[string]any{}}
	jsonLD := findJSONLD(doc)
	if value, ok := jsonLD["name"].(string); ok && value != "" {
		offer.Name = &value
	}
	if value, ok := jsonLD["address"].(string); ok && value != "" {
		offer.Address = &value
	}
	if object, ok := jsonLD["address"].(map[string]any); ok {
		if value, ok := object["streetAddress"].(string); ok && value != "" {
			offer.Address = &value
		}
	}
	if object, ok := jsonLD["offers"].(map[string]any); ok {
		if value, ok := numberFromJSON(object["price"]); ok {
			offer.RentYen = &value
		}
	}
	if value, ok := jsonLD["floorSize"].(map[string]any); ok {
		if area, ok := floatFromJSON(value["value"]); ok {
			offer.AreaSquareMeters = &area
		}
	}
	if value, ok := jsonLD["numberOfRooms"].(string); ok && value != "" {
		offer.RoomLayout = &value
	}
	for key, value := range jsonLD {
		offer.Details[key] = value
	}
	return offer
}

func findJSONLD(doc *html.Node) map[string]any {
	var result map[string]any
	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if result != nil {
			return
		}
		if node.Type == html.ElementNode && node.Data == "script" {
			for _, attr := range node.Attr {
				if attr.Key != "type" || attr.Val != "application/ld+json" {
					continue
				}
				var text strings.Builder
				for child := node.FirstChild; child != nil; child = child.NextSibling {
					if child.Type == html.TextNode {
						text.WriteString(child.Data)
					}
				}
				var value map[string]any
				if json.Unmarshal([]byte(strings.TrimSpace(text.String())), &value) == nil {
					result = value
				}
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(doc)
	if result == nil {
		return map[string]any{}
	}
	return result
}

func numberFromJSON(value any) (int64, bool) {
	switch value := value.(type) {
	case float64:
		return int64(value), true
	case string:
		value = strings.ReplaceAll(value, ",", "")
		parsed, err := strconv.ParseInt(value, 10, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func floatFromJSON(value any) (float64, bool) {
	switch value := value.(type) {
	case float64:
		return value, true
	case string:
		value = strings.ReplaceAll(value, ",", "")
		parsed, err := strconv.ParseFloat(value, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}
