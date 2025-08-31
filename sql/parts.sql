CREATE TABLE components (
    id SERIAL PRIMARY KEY,
    added_by INT NOT NULL,
    name VARCHAR(64),
    description VARCHAR(255),
    footprint VARCHAR(64),
    manufacturer VARCHAR(64),
    supplier VARCHAR(64),
    amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY(added_by) REFERENCES users(id)
);

-- add random components
INSERT INTO components (added_by, name, description, footprint, manufacturer, supplier, amount)
SELECT
    2 AS added_by,
    'Component_' || gs AS name,
    'Random electronic part ' || gs || ' for testing' AS description,
    (ARRAY['SOT-23', 'QFN-32', 'DIP-8', 'SOIC-16', 'TSSOP-20', 'BGA-100'])[floor(random() * 6 + 1)] AS footprint,
    (ARRAY['Texas Instruments', 'Analog Devices', 'NXP', 'Microchip', 'STMicroelectronics', 'Infineon'])[floor(random() * 6 + 1)] AS manufacturer,
    (ARRAY['Digikey', 'Mouser', 'Arrow', 'RS Components', 'Element14'])[floor(random() * 5 + 1)] AS supplier,
    floor(random() * 100 + 1)::INT AS amount
FROM generate_series(1, 10000) gs;
