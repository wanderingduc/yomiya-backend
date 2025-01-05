CREATE DATABASE test_db;
USE test_db;

CREATE TABLE users(
    username VARCHAR(255) PRIMARY KEY,
    passwd VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE reports(
    report_id SERIAL,
    user_fk VARCHAR(255),
    bug TEXT,
    FOREIGN KEY (user_fk) REFERENCES users(username)
);

CREATE TABLE authors(
    author_id varchar(255) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL
);

CREATE TABLE books(
    -- book_id VARCHAR(13) PRIMARY KEY,
    book_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author_fk VARCHAR(255),
    summary TEXT, -- CHANGE NOT MADE TO DB
    FOREIGN KEY (author_fk) REFERENCES authors(author_id)
);

CREATE TABLE libs(
    lib_id VARCHAR(255) PRIMARY KEY,
    lib_name VARCHAR(255) NOT NULL,
    user_fk VARCHAR(255),
    FOREIGN KEY (user_fk) REFERENCES users(username)
);

CREATE TABLE lib(
    lib_fk VARCHAR(255) NOT NULL,
    -- book_fk VARCHAR(13) NOT NULL,
    book_fk SERIAL,
    FOREIGN KEY (lib_fk) REFERENCES libs(lib_id),
    FOREIGN KEY (book_fk) REFERENCES books(book_id),
    PRIMARY KEY (lib_fk, book_fk)
);

CREATE TABLE new_books(
    contrib_id BINARY(16) PRIMARY KEY,
    book_id VARCHAR(13),
    title VARCHAR(255),
    author VARCHAR(255),
    user_fk VARCHAR(255),
    FOREIGN KEY (user_fk) REFERENCES users(username)
);

ALTER TABLE books ADD FULLTEXT(title, author_fk);
ALTER TABLE libs ADD FULLTEXT(lib_id, lib_name);
ALTER TABLE new_books ADD FULLTEXT(title, author);
