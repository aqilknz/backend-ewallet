CREATE TABLE IF NOT EXISTS profiles (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users(id),
    full_name VARCHAR(100),
    phone VARCHAR(20),
    photo TEXT DEFAULT 'https://i.pravatar.cc/150?img=0',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP 
);