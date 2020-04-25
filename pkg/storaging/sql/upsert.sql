INSERT INTO localeitems ( key, bundle, lang, content ) 
VALUES( $1,$2,$3,$4)
ON CONFLICT ON CONSTRAINT ukey_localeitems
DO UPDATE SET content = $4 
WHERE localeitems.key = $1 AND localeitems.bundle = $2 AND localeitems.lang = $3
RETURNING id;