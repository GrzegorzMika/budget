CREATE TABLE IF NOT EXISTS categories (
    name TEXT PRIMARY KEY
);

INSERT INTO categories (name) VALUES 
    ('Środki czystości'),
    ('Spożywcze'),
    ('Odzież i obuwie'),
    ('Rozrywka'),
    ('Gastronomia'),
    ('Kosmetyki'),
    ('Edukacja'),
    ('Transport'),
    ('Zdrowie'),
    ('Wyposażenie'),
    ('Wypoczynek i podróże'),
    ('Czynsz i media'),
    ('Elektronika'),
    ('Prezenty'),
    ('Ubezpieczenia'),
    ('Samochód'),
    ('Ślub'),
    ('Artykuły higieniczne'),
    ('Inne usługi'),
    ('Budowa domu')
ON CONFLICT DO NOTHING;
