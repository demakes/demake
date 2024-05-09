CREATE TABLE demake_version (
    version_num integer NOT NULL
);

INSERT INTO demake_version (version_num) VALUES (1);

{{$sqlite:=false}}

{{if eq .DBType "sqlite3"}}
    {{$sqlite = true}}
{{end}}

/* Nodes */

CREATE TABLE node (
    {{if $sqlite}}
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    {{else}}
    id BIGINT NOT NULL,
    {{end}}
    hash bytea NOT NULL,
    type character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    data bytea
);


{{ if not $sqlite }}
CREATE SEQUENCE node_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE node_seq OWNED BY node.id;

ALTER TABLE ONLY node ALTER COLUMN id SET DEFAULT nextval('node_seq'::regclass);
ALTER TABLE ONLY node
    ADD CONSTRAINT node_pkey PRIMARY KEY (id);

{{ end }}

CREATE UNIQUE INDEX ix_node_hash ON node (hash) WHERE (deleted_at IS NULL);
CREATE INDEX ix_node_created_at ON node (created_at);
CREATE INDEX ix_node_deleted_at ON node (deleted_at);
CREATE INDEX ix_node_updated_at ON node (updated_at);
CREATE INDEX ix_node_type ON node (type);

/* Edges */

CREATE TABLE edge (
    {{if $sqlite}}
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    {{else}}
    id BIGINT NOT NULL,
    {{end}}
    ext_id bytea NOT NULL,
    {{ if $sqlite }}
    from_id INTEGER NOT NULL REFERENCES node(id),   
    to_id INTEGER NOT NULL REFERENCES node(id),   
    {{else}}
    from_id bigint NOT NULL REFERENCES node(id),
    to_id bigint NOT NULL REFERENCES node(id),
    {{end}}
    name character varying NOT NULL,
    key character varying, -- only defined for map-based edges
    ind integer,
    type integer,
    follow bool DEFAULT TRUE,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone,
    data bytea
);

{{ if not $sqlite}}

CREATE SEQUENCE edge_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE edge_seq OWNED BY edge.id;
ALTER TABLE ONLY edge ALTER COLUMN id SET DEFAULT nextval('edge_seq'::regclass);

ALTER TABLE ONLY edge
    ADD CONSTRAINT edge_pkey PRIMARY KEY (id);

{{ end }}

CREATE INDEX ix_edge_created_at ON edge (created_at);
CREATE INDEX ix_edge_deleted_at ON edge (deleted_at);
CREATE INDEX ix_edge_follow ON edge (follow);
CREATE UNIQUE INDEX ix_edge_ext_id ON edge (ext_id) WHERE (deleted_at IS NULL);
CREATE UNIQUE INDEX ix_edge_unique ON edge (from_id, to_id, name, ind, key, type) WHERE (deleted_at IS NULL);
CREATE INDEX ix_edge_name ON edge (name);
CREATE INDEX ix_edge_key ON edge (key);
CREATE INDEX ix_edge_ind ON edge (ind);
CREATE INDEX ix_edge_type ON edge (type);
CREATE INDEX ix_edge_updated_at ON edge (updated_at);
CREATE INDEX ix_edge_from_to_id ON edge (from_id, to_id);
