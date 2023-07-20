CREATE TABLE "category" (
    "id" UUID PRIMARY KEY,
    "title" VARCHAR NOT NULL,
    "parent_id" UUID REFERENCES "category"("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);


CREATE TABLE "product" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "category_id" UUID REFERENCES "category"("id"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);

CREATE TABLE "users" (
    "id" UUID PRIMARY KEY,
    "first_name" VARCHAR NOT NULL,
    "last_name" VARCHAR NOT NULL,
    "balans"  NUMERIC DEFAULT 0,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);
SELECT
			p.id,
			p.name,
			p.price,
			c.id,         
			c.title,      
			c.created_at, 	
			c.updated_at,
			p.created_at,
			p.updated_at
		
		FROM product as p
        join category as c on c.id=p.category_id
		WHERE p.id = '0b0d4901-9703-4ada-958a-887d274bda99';
SELECT
    name,
    price

from product as p
join category as c on c.id=p.category_id

where p.id ='0b0d4901-9703-4ada-958a-887d274bda99'
;

