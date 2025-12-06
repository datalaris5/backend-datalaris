CREATE TABLE tenants (
	id bigserial PRIMARY KEY,
	"name" varchar(255) NOT NULL,
	description text NULL,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	tenant_key varchar(255) NULL,
	CONSTRAINT tenants_name_key UNIQUE (name)
);

CREATE TABLE roles (
	id bigserial PRIMARY KEY,
	"name" varchar(255) NOT NULL,
	tenant_id BIGINT REFERENCES tenants(id) ON DELETE CASCADE,   
	role_type varchar(50) NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	is_active bool DEFAULT true NULL
);

CREATE TABLE users (
	id bigserial PRIMARY KEY,
	"name" varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	"password" varchar(255) NOT NULL,
	tenant_id BIGINT REFERENCES tenants(id) ON DELETE CASCADE,
	role_id BIGINT REFERENCES roles(id) ON DELETE CASCADE,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	created_by int8 NULL,
	updated_by int8 NULL,
	is_active bool DEFAULT true NULL
);

CREATE TABLE shopee_data_upload_details (
	id bigserial PRIMARY KEY,
	tenant_id BIGINT REFERENCES tenants(id) ON DELETE CASCADE,    
	tanggal date,
	total_penjualan NUMERIC(18,2),
	total_pesanan int,
	penjualan_per_pesanan NUMERIC(18, 2),
	produk_klik int,
	total_pengunjung int,
	tingkat_konversi_harian NUMERIC(5,2),
	pesanan_dibatalkan int,
	pesanan_dikembalikan int,
	penjualan_dikembalikan int,
	pembeli int,
	total_pembeli_baru int,
	total_pembeli_saat_ini int,
	total_potensi_pembeli int,
	tingkat_pembelian_berulang NUMERIC(5,2),
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	created_by int8 NULL,
	updated_by int8 NULL
);

CREATE TABLE shopee_data_upload_summaries (
	id bigserial PRIMARY KEY,
	tenant_id BIGINT REFERENCES tenants(id) ON DELETE CASCADE,    
	tanggal date,
	total_penjualan NUMERIC(18,2),
	total_pesanan int,
	penjualan_per_pesanan NUMERIC(18, 2),
	produk_klik int,
	total_pengunjung int,
	tingkat_konversi_harian NUMERIC(5,2),
	pesanan_dibatalkan int,
	pesanan_dikembalikan int,
	penjualan_dikembalikan int,
	pembeli int,
	total_pembeli_baru int,
	total_pembeli_saat_ini int,
	total_potensi_pembeli int,
	tingkat_pembelian_berulang NUMERIC(5,2),
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	created_by int8 NULL,
	updated_by int8 NULL
);

CREATE TABLE stores (
	id bigserial PRIMARY KEY,
	tenant_id BIGINT REFERENCES tenants(id) ON DELETE CASCADE,
	marketplace_id BIGINT REFERENCES marketplaces(id) ON DELETE CASCADE,
	"name" varchar(255) NOT NULL,
	is_active bool DEFAULT true,
	is_deleted bool DEFAULT false,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	created_by int8 NULL,
	updated_by int8 NULL
);

CREATE TABLE marketplaces (
	id bigserial PRIMARY KEY,
	"name" varchar(255) NOT NULL,
	is_active bool DEFAULT true,
	is_deleted bool DEFAULT false,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	created_by int8 NULL,
	updated_by int8 NULL
);




