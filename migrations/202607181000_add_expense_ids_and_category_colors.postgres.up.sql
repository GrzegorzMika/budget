ALTER TABLE expenses ADD COLUMN IF NOT EXISTS id BIGSERIAL PRIMARY KEY;

ALTER TABLE categories ADD COLUMN IF NOT EXISTS color TEXT NOT NULL DEFAULT 'c1';

-- Spread existing categories across the 16-color palette deterministically
UPDATE categories SET color = 'c' || (1 + (abs(hashtext(name)) % 16));
