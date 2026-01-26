-- Veritabanı oluştur.
CREATE DATABASE lesson_db;

-- Veritabanı kopyalama.
CREATE DATABASE lesson_db_copy
WITH TEMPLATE lesson_db;

-- Veritabanı adıni değiştir.
ALTER DATABASE lesson_db
RENAME TO my_lesson_db;

-- Eğer varsa veritabanını sil.
DROP DATABASE IF EXISTS my_lesson_db;

-- users tablosu yoksa oluştur.
CREATE TABLE IF NOT EXISTS users (
    -- id SERIAL PRIMARY KEY, -- SERIAL HACKING yaparak otomatik artan birincil anahtar oluştur.
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, -- Yukardaki ile aynı işi yapar ama daha standart bir yöntemdir.
    -- username VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL CHECK (LENGTH(username) >= 3), -- username en az 3 karakter olmalı.
    email VARCHAR(100) UNIQUE,
    password VARCHAR(100) NOT NULL, -- Burada şifre hashlenmiş olarak saklanacak.
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- kayıt oluşturulma zamanında şimdiki zamanı alır.
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- kayıt güncellenme zamanında şimdiki zamanı alır. Fakat bunu tetiklemek için bir TRIGGER yazmak gerekir.
);
-- TIMESTAMPTZ kullanılırsa zaman dilimi bilgisi de saklanır.

-- Tablo adını değiştirme
ALTER TABLE users
RENAME TO app_users;

-- Tablo kopyalama
CREATE TABLE app_users_backup AS TABLE app_users;

-- Tabloyu boşaltma (truncate)
TRUNCATE TABLE users;

-- users tablosu varsa sil.
DROP TABLE IF EXISTS app_users;


-- Trigger fonksiyonu oluşturma.
CREATE
OR REPLACE FUNCTION update_updated_at () RETURNS TRIGGER AS $$
BEGIN
-- Güncelleme işlemi sırasında updated_at sütununu şimdiki zamanla güncelle. Ve NEW kaydını döndür.
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger fonksiyonu silme, syntax: DROP FUNCTION IF EXISTS functionname();
DROP FUNCTION IF EXISTS update_updated_at();

-- Trigger oluşturma ve tabloya bağlayip fonksiyonu tetikleme.
CREATE TRIGGER trg_users_updated
BEFORE UPDATE ON users
FOR EACH ROW -- Her satır için tetikle. Burasi ROW veya STATEMENT olabilir. STATEMENT tüm tablo için tetikler, ROW ise her satır için tetikler.
EXECUTE FUNCTION update_updated_at();

-- Trigger silme, syntax: DROP TRIGGER IF EXISTS triggername ON tablename;
DROP TRIGGER IF EXISTS trg_users_updated ON users;

-- ALTER Kullanımı: Genel olarak tablo yapısını değiştirmek için kullanılır.
-- Sütun ekleme
ALTER TABLE users
ADD COLUMN last_login TIMESTAMP;

-- Sütun adını değiştirme
ALTER TABLE users
RENAME COLUMN is_active TO status;

-- Sütun defalt değerini değiştirme
ALTER TABLE users
ALTER COLUMN status SET DEFAULT FALSE;

-- Sütun veritipini değiştirme
ALTER TABLE users
ALTER COLUMN last_login TYPE TIMESTAMPTZ;

ALTER TABLE users
ALTER COLUMN status TYPE VARCHAR(20); -- status sütununu VARCHAR(20) tipine çevirdik.

ALTER TABLE users
ALTER COLUMN status TYPE BOOLEAN USING CASE
	WHEN status = 'true' THEN TRUE
	WHEN status = 'false' THEN FALSE
	ELSE NULL
END; -- Bu şekilde string ifadeleri boolean'a çevirebiliriz.

ALTER TABLE users
ALTER COLUMN status TYPE BOOLEAN USING status::boolean; -- Bu şekilde de string ifadeleri boolean'a çevirebiliriz. Ama CASE yapısı daha esnek cunku farklı string ifadeleri de yakalayabiliriz.

-- Yeni bir is_active sütunu ekledik. Default değeri 'active'. (sadece passive/active gibi değerler alabilir.)
-- Eger is_active sütununu INT tipine çevirmek istersek ve active/passive yerine 1/0 gibi değerler almak istersek:
-- Once DROP DEFAULT ile default değeri kaldıralım.
-- Sonra mevcut string değerleri INT'e çevirelim ve using ile mevcut kayıtlarda bu dönüşümü yapalım. Kayitlarda active ise 1, inactive ise 0 yapalım farklı bir değer varsa NULL yapalım.
-- Son olarak tekrar default değeri set edelim.
ALTER TABLE users
ADD COLUMN is_active VARCHAR(20) DEFAULT 'active';

ALTER TABLE users
ALTER COLUMN is_active DROP DEFAULT;

ALTER TABLE users
ALTER COLUMN is_active TYPE INT USING CASE
	WHEN is_active = 'active' THEN 1
	WHEN is_active = 'inactive' THEN 0
	ELSE NULL
END;

ALTER TABLE users
ALTER COLUMN is_active SET DEFAULT 1;

-- Sütun silme
ALTER TABLE users
DROP COLUMN IF EXISTS last_login;

-- Tabloya  açıklama ekleme
COMMENT ON TABLE users IS 'Kullanıcı bilgilerini tutan tablo';
-- Sütuna açıklama ekleme
COMMENT ON COLUMN users.username IS 'Kullanıcının username fieldi benzersiz olmalıdır';
-- Tablo açıklamasını silme
COMMENT ON TABLE users IS NULL;
-- Sütun açıklamasını silme
COMMENT ON COLUMN users.email IS NULL;



-- BLOG ICIN DB ISLEMLERI
-- categories tablosu oluşturma
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(60) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- articles tablosu oluşturma
CREATE TABLE IF NOT EXISTS articles (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    slug VARCHAR(60) UNIQUE NOT NULL,
    short_description VARCHAR(150),
    description TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    user_id BIGINT REFERENCES users (id) ON DELETE SET NULL, -- ON DELETE SET NULL: Eğer kullanıcı silinirse, user_id NULL olur. ON UPDATE CASCADE : Eğer kullanıcı id'si değişirse, burada da otomatik güncellenir.
    seo_settings JSONB, -- JSONB tipinde SEO ayarlarını saklayabiliriz. Örnek: {"meta_title": "Başlık", "meta_description": "Açıklama", "keywords": ["anahtar1", "anahtar2"]}
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- articles_categories ilişki tablosu oluşturma (Many-to-Many ilişkisi için)
CREATE TABLE IF NOT EXISTS article_categories (
    article_id BIGINT NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories (id) ON DELETE CASCADE,
    PRIMARY KEY (article_id, category_id) -- Bir makale bir kategoride sadece bir kez bulunabilir.
);

-- index oluşturma, syntax: CREATE INDEX idx_tablename_columnname ON tablename(columnname);
CREATE INDEX IF NOT EXISTS idx_articles_user_id ON articles (user_id); -- Kullanıcıya göre makaleleri hızlıca bulmak için index.
CREATE INDEX IF NOT EXISTS idx_article_categories_category_id ON article_categories (category_id); -- Kategoriye gore makaleleri hızlıca bulmak için index.
-- CREATE INDEX IF NOT EXISTS idx_article_categories_article_id ON article_categories (article_id); -- Makaleye gore kategorileri hızlıca bulmak için index.

-- Index silme, syntax: DROP INDEX IF EXISTS indexname;
-- DROP INDEX IF EXISTS idx_articles_user_id;
-- DROP INDEX IF EXISTS idx_article_categories_category_id;
-- DROP INDEX IF EXISTS idx_article_categories_article_id;


-- INSERT INTO users tablosuna örnek veriler ekleme
INSERT INTO
    users (username, email, password)
VALUES
    ('sercanozen', 'sercan@example.com', 'hashed_password_1'),
    ('alinozen', 'alin@example.com', 'hashed_password_2');

-- Egerki conflict (çakışma) durumu olursa ne yapacağını belirtme.
-- Örneğin email alanı UNIQUE olduğu için aynı email ile tekrar ekleme yapmaya çalışırsak hata alırız.
-- Bu durumda ON CONFLICT ifadesi ile ne yapacağını belirtebiliriz.
-- Örneğin email alanında çakışma olursa, username ve password alanlarını güncelle.
INSERT INTO
    users (username, email, password)
VALUES
    ('new_sercanozen', 'sercan1@example.com', 'new_hashed_password_1'),
    ('new_alinozen', 'alin1@example.com', 'new_hashed_password_2')
ON CONFLICT (email) DO
UPDATE
    SET
        username = EXCLUDED.username,
        password = EXCLUDED.password;
-- EXCLUDED ifadesi, çakışan (conflicting) satırdaki değerleri temsil eder.
-- Eğer email zaten varsa, hiçbir şey yapma.

-- INSERT INTO articles tablosuna örnek veriler ekleme
INSERT INTO
    articles (title, slug, short_description, description, is_active, user_id, seo_settings)
VALUES
    (
        'PostgreSQL ve Golang Kurulumu',
        'postgresql-ve-golang-kurulumu',
        'Bu makalede PostgreSQL ve Golang kurulumu anlatılacaktır.',
        '<p>PostgreSQL kurulumu için öncelikle ...</p>',
        TRUE,
        15,
        '{"meta_title": "PostgreSQL ve Golang Kurulumu", "meta_description": "Bu makalede PostgreSQL ve Golang kurulumu anlatılacaktır.", "keywords": ["postgresql", "golang", "kurulum"]}'
    ),
    (
        'Go ve PostgreSQL Rehberi',
        'go-ve-postgresql-rehberi',
        'Go programlama dili ile PostgreSQL veritabanı kullanımı hakkında rehber.',
        '<p>Go ile PostgreSQL kullanmak için öncelikle ...</p>',
        FALSE,
        12,
        '{"meta_title": "Go ve PostgreSQL Rehberi", "meta_description": "Go programlama dili ile PostgreSQL veritabanı kullanımı hakkında rehber.", "keywords": ["go", "postgresql", "rehber"]}'
    );

-- INSERT INTO categories tablosuna örnek veriler ekleme
INSERT INTO
    categories (name, slug)
VALUES
    ('Go Programlama', 'go-programlama'),
    ('Veritabanı', 'veritabani'),
    ('Web Geliştirme', 'web-gelistirme');

-- INSERT INTO article_categories ilişki tablosuna örnek veriler ekleme
INSERT INTO
    article_categories (article_id, category_id)
VALUES
    (7, 1), -- İlk makale Go Programlama kategorisinde
    (7, 2), -- İlk makale Veritabanı kategorisinde
    (8, 1), -- İkinci makale Go Programlama kategorisinde
    (8, 3); -- İkinci makale Web Geliştirme kategorisinde