CREATE DATABASE test_db;
USE test_db;

CREATE TABLE users(
    username VARCHAR(255) PRIMARY KEY,
    passwd VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE authors(
    author_id varchar(255) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL
);

CREATE TABLE books(
    book_id VARCHAR(13) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author_fk VARCHAR(255),
    FOREIGN KEY (author_fk) REFERENCES authors(author_id)
);

CREATE TABLE libs(
    lib_id VARCHAR(255) PRIMARY KEY,
    lib_name VARCHAR(255) NOT NULL,
    user_fk VARCHAR(255),
    book_fk VARCHAR(13),
    FOREIGN KEY (user_fk) REFERENCES users(username),
    FOREIGN KEY (book_fk) REFERENCES books(book_id)
);

CREATE TABLE lib(
    lib_fk VARCHAR(255) NOT NULL,
    book_fk VARCHAR(13) NOT NULL,
    FOREIGN KEY (lib_fk) REFERENCES libs(lib_id),
    FOREIGN KEY (book_fk) REFERENCES books(book_id)
);

CREATE TABLE new_books(
    contrib_id BINARY(16) PRIMARY KEY,
    book_id VARCHAR(13),
    title VARCHAR(255),
    author VARCHAR(255),
    user_fk VARCHAR(255),
    FOREIGN KEY (user_fk) REFERENCES users(username)
);
