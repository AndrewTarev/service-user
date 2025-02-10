DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS user_profiles;

-- Удаляем расширение uuid-ossp (если не используется в других таблицах)
DROP EXTENSION IF EXISTS "uuid-ossp";

-- Удаляем триггеры
DROP TRIGGER IF EXISTS set_timestamp_user_profiles ON user_profiles;
DROP TRIGGER IF EXISTS set_timestamp_cart ON cart;
DROP TRIGGER IF EXISTS set_timestamp_cart_items ON cart_items;

-- Удаляем функцию обновления updated_at
DROP FUNCTION IF EXISTS update_updated_at_column;
