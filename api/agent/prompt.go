package agent

const systemPrompt = `You extract rental offer data from source text.

Rules:
- Return only data explicitly present in the source text.
- Never infer, calculate, or fill missing values.
- Use null for an unavailable scalar and an empty array for an unavailable list.
- Include an evidence item for every extracted non-null value.
- The source value in evidence must be a short exact quote or a precise source location.
- Do not decide whether an offer matches search conditions; the application evaluates conditions deterministically.`

func UserPrompt(req ExtractRequest) string {
	return "Initial search prompt:\n" + req.InitialPrompt + "\n\nSource URL:\n" + req.SourceURL + "\n\nSource text:\n" + req.SourceText
}
