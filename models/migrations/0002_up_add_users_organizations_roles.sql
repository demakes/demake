UPDATE demake_version SET version_num = 2;

{{$sqlite:=false}}

{{if eq .DBType "sqlite3"}}
    {{$sqlite = true}}
{{end}}


CREATE TABLE organization (
    {{if $sqlite}}
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    {{else}}
    id BIGINT NOT NULL,
    {{end}}
    ext_id bytea NOT NULL,
    source_id bytea NOT NULL,
    source character varying NOT NULL,
    name character varying NOT NULL,
    description character varying DEFAULT '' NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    data jsonb
);


{{ if not $sqlite }}
CREATE SEQUENCE organization_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE organization_seq OWNED BY organization.id;

ALTER TABLE ONLY organization ALTER COLUMN id SET DEFAULT nextval('organization_seq'::regclass);
ALTER TABLE ONLY organization
    ADD CONSTRAINT organization_pkey PRIMARY KEY (id);

{{ end }}

CREATE UNIQUE INDEX ix_organization_source_source_id ON organization (source, source_id);
CREATE UNIQUE INDEX ix_organization_ext_id ON organization (ext_id);
CREATE INDEX ix_organization_created_at ON organization (created_at);
CREATE INDEX ix_organization_deleted_at ON organization (deleted_at);
CREATE INDEX ix_organization_updated_at ON organization (updated_at);
CREATE INDEX ix_organization_source ON organization (source, source_id);

/* Users */

CREATE TABLE "user" (
    {{if $sqlite}}
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    {{else}}
    id BIGINT NOT NULL,
    {{end}}
    ext_id bytea NOT NULL,
    display_name character varying NOT NULL,
    source character varying NOT NULL,
    source_id bytea NOT NULL,
    email character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    data jsonb
);

{{ if not $sqlite}}

CREATE SEQUENCE user_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE user_seq OWNED BY "user".id;
ALTER TABLE ONLY "user" ALTER COLUMN id SET DEFAULT nextval('user_seq'::regclass);

ALTER TABLE ONLY "user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);

{{ end }}

CREATE INDEX ix_user_created_at ON "user" (created_at);
CREATE INDEX ix_user_updated_at ON "user" (updated_at);
CREATE INDEX ix_user_deleted_at ON "user" (deleted_at);
CREATE UNIQUE INDEX ix_user_ext_id ON "user" (ext_id);
CREATE UNIQUE INDEX ix_user_email ON "user" (email);
CREATE INDEX ix_user_display_name ON "user" (display_name);
CREATE INDEX ix_user_source ON "user" (source, source_id);
CREATE INDEX user_id ON "user" (id);

/* User roles */

CREATE TABLE user_role (
    {{if $sqlite}}
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    {{else}}
    id BIGINT NOT NULL,
    {{end}}
    ext_id bytea NOT NULL,
    {{ if $sqlite }}
    organization_id INTEGER NOT NULL REFERENCES organization(id),   
    {{else}}
    organization_id bigint NOT NULL REFERENCES organization(id),
    {{end}}
    {{ if $sqlite }}
    user_id INTEGER NOT NULL REFERENCES "user"(id),   
    {{else}}
    user_id bigint NOT NULL REFERENCES "user"(id),
    {{end}}
    role CHARACTER VARYING NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp without time zone,
    data jsonb
);

{{ if not $sqlite}}

CREATE SEQUENCE user_role_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE user_role_seq OWNED BY user_role.id;
ALTER TABLE ONLY user_role ALTER COLUMN id SET DEFAULT nextval('user_role_seq'::regclass);

ALTER TABLE ONLY user_role
    ADD CONSTRAINT user_role_pkey PRIMARY KEY (id);

{{ end }}

CREATE INDEX ix_user_role_created_at ON user_role (created_at);
CREATE INDEX ix_user_role_updated_at ON user_role (updated_at);
CREATE INDEX ix_user_role_deleted_at ON user_role (deleted_at);
CREATE UNIQUE INDEX ix_user_role_ext_id ON user_role (ext_id);
CREATE INDEX ix_user_role_display_role ON user_role (role);
CREATE INDEX user_role_id ON user_role (id);
