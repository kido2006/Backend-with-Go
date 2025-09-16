-- 1. Thêm cột role_id với default = 1
ALTER TABLE users
  ADD COLUMN role_id BIGINT DEFAULT 1;

-- 2. Thêm khóa ngoại (nếu muốn ràng buộc với bảng roles)
ALTER TABLE users
  ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id);

-- 3. Update tất cả user hiện tại thành role 'user'
UPDATE users
SET role_id = (
  SELECT id FROM roles WHERE name = 'user' LIMIT 1
);

-- 4. Xóa default
ALTER TABLE users
  ALTER role_id DROP DEFAULT;

-- 5. Bắt buộc role_id phải NOT NULL
ALTER TABLE users
  MODIFY role_id BIGINT NOT NULL;
