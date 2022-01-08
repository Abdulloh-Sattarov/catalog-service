create table categories (
    category_id uuid primary key,
    name varchar(64),
    parent_uuid uuid default null,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp default null
);

create table authors (
     author_id uuid primary key,
     name varchar(64),
     created_at timestamp default current_timestamp,
     updated_at timestamp default current_timestamp,
     deleted_at timestamp default null
);

create table books (
       book_id uuid primary key,
       name varchar(64),
       author_id uuid references authors(author_id),
       price decimal default null,
       created_at timestamp default current_timestamp,
       updated_at timestamp default current_timestamp,
       deleted_at timestamp default null
);

create table book_categories (
     book_category_id serial primary key,
     book_id uuid references books(book_id),
     category_id uuid references categories(category_id)
);