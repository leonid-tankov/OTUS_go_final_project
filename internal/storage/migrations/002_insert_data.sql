-- +goose Up
INSERT INTO slots (description)
VALUES
    ('main'),
    ('left'),
    ('right'),
    ('bottom') ON CONFLICT DO NOTHING;

INSERT INTO social_dem_groups (description)
VALUES
    ('men 18-30'),
    ('men 31-60'),
    ('men >60'),
    ('women 18-30'),
    ('women 31-60'),
    ('women >60') ON CONFLICT DO NOTHING;

INSERT INTO banners (description)
VALUES
    ('yandex'),
    ('mail'),
    ('google'),
    ('youtube'),
    ('facebook'),
    ('amazon'),
    ('netflix'),
    ('yahoo'),
    ('reddit'),
    ('dzen'),
    ('vk'),
    ('zoom'),
    ('pinterest'),
    ('wikipedia'),
    ('linkedin'),
    ('samsung'),
    ('microsoft'),
    ('ebay') ON CONFLICT DO NOTHING;

