CREATE TABLE public.links (
	id uuid NOT NULL,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	deleted_at timestamptz NULL,
	shortlink varchar NOT NULL,
	longlink varchar NOT NULL,
	"data" varchar NULL,	
	CONSTRAINT links_pk PRIMARY KEY (id)
);