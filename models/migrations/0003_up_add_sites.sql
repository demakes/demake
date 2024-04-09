UPDATE klaro_version SET version_num = 3;

{{$sqlite:=false}}

{{if eq .DBType "sqlite3"}}
    {{$sqlite = true}}
{{end}}

/* Sites */

CREATE TABLE site (
    {{if $sqlite}}
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    {{else}}
    id BIGINT NOT NULL,
    {{end}}
    ext_id bytea NOT NULL,
    description character varying DEFAULT '' NOT NULL,
    name character varying NOT NULL,
    {{ if $sqlite }}
    head_id INTEGER REFERENCES node(id),
    {{else}}
    head_id bigint REFERENCES node(id),
    {{end}}
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    data jsonb
);

{{ if not $sqlite}}

CREATE SEQUENCE site_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE site_seq OWNED BY site.id;
ALTER TABLE ONLY site ALTER COLUMN id SET DEFAULT nextval('site_seq'::regclass);

ALTER TABLE ONLY site
    ADD CONSTRAINT site_pkey PRIMARY KEY (id);

{{ end }}

CREATE INDEX ix_site_created_at ON site (created_at);
CREATE INDEX ix_site_deleted_at ON site (deleted_at);
CREATE UNIQUE INDEX ix_site_ext_id ON site (ext_id);
CREATE INDEX ix_site_name ON site (name);
CREATE INDEX ix_site_updated_at ON site (updated_at);
CREATE INDEX site_id ON site (id);

