ALTER TABLE properties RENAME TO rental_offers;

ALTER TABLE rental_offers
    RENAME COLUMN source_property_id TO source_offer_id;

ALTER TABLE rental_offers
    RENAME CONSTRAINT properties_search_profile_id_fkey TO rental_offers_search_profile_id_fkey;

ALTER INDEX properties_profile_source_property_id_idx
    RENAME TO rental_offers_profile_source_offer_id_idx;

ALTER INDEX properties_search_profile_id_idx
    RENAME TO rental_offers_search_profile_id_idx;
