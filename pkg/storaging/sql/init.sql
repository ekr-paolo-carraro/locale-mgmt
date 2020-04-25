CREATE TABLE IF NOT EXISTS localeitems(
    id serial NOT NULL,
    key VARCHAR(512),
    bundle VARCHAR(128),
    lang VARCHAR(8),
    content VARCHAR(4096),
    CONSTRAINT 
        pKey_localeitems PRIMARY KEY (id),
	CONSTRAINT
        uKey_localeitems UNIQUE ( key, bundle, lang ) 
)