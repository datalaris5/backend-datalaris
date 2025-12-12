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
	store_id BIGINT REFERENCES stores(id) ON DELETE CASCADE,    
	tanggal date,
	total_penjualan NUMERIC(18,2),
	total_pesanan int,
	penjualan_per_pesanan NUMERIC(18, 2),
	produk_klik int,
	total_pengunjung int,
	tingkat_konversi_harian NUMERIC(5,2),
	pesanan_dibatalkan int,
	penjualan_dibatalkan int,
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

CREATE TABLE shopee_data_upload_iklan_details (
	id bigserial PRIMARY KEY,
	store_id BIGINT REFERENCES stores(id) ON DELETE CASCADE,  
	tanggal DATE,  
	nama_iklan TEXT,
	status VARCHAR(200),
	jenis_iklan VARCHAR(100),
	kode_produk VARCHAR(100),
	tampilan_iklan TEXT,
	mode_bidding VARCHAR(100),
	penempatan_iklan VARCHAR(150),
	tanggal_mulai VARCHAR(100),
	tanggal_selesai VARCHAR(100),
	dilihat INT,
	jumlah_klik INT,
	presentase_klik NUMERIC(18,2),
	konversi INT,
	konversi_langsung INT,
	tingkat_konversi NUMERIC(18,2),
	biaya_per_konversi NUMERIC(18,2),
	biaya_per_konversi_langsung NUMERIC(18,2),
	produk_terjual INT,
	terjual_langsung INT,
	omzet_penjualan NUMERIC(18,2),
	gvm_langsung NUMERIC(18,2),
	biaya NUMERIC(18,2),
	efektifitas_iklan NUMERIC(18,2),
	efektifitas_langsung NUMERIC(18,2),
	acos NUMERIC(18,2),
	acos_langsung NUMERIC(18,2),
	jumlah_produk_dilihat INT,
	jumlah_produk_diklik INT,
	presentase_produk_diklik NUMERIC(18,2),
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	created_by int8 NULL,
	updated_by int8 NULL
);

CREATE TABLE shopee_data_upload_chat_details (
	id bigserial PRIMARY KEY,
	store_id BIGINT REFERENCES stores(id) ON DELETE CASCADE,  
	tanggal DATE,  
	pengunjung INT,
	jumlah_chat INT,
	pengunjung_bertanya INT,
	pertanyaan_diajukan NUMERIC(18,2),
	chat_dibalas INT,
	chat_belum_dibalas INT,
	waktu_respon_rata_rata INTERVAL,
	csat NUMERIC(18,2),
	waktu_respon_chat_pertama INTERVAL,
	presentase_chat_dibalas NUMERIC(18,2),
	tingkat_konversi_jumlah_chat_direspon NUMERIC(18,2),
	total_pembeli INT,
	total_pesanan INT,
	produk INT,
	penjualan NUMERIC(18,2),
	tingkat_konversi_chat_dibalas NUMERIC(18,2),
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	created_by int8 NULL,
	updated_by int8 NULL
);

CREATE TABLE history_data_uploads (
	id bigserial PRIMARY KEY,
	store_id BIGINT REFERENCES stores(id) ON DELETE CASCADE,    
	filename VARCHAR(250),
	status VARCHAR(20),
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




