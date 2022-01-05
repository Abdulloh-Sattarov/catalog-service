drop table if exists categories;
drop table if exists authors;
drop table if exists books;
drop table if exists book_categories;

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

insert into categories(category_id, name) values ('2bd05294-e4c5-46ba-a458-ba54c799e4e3', 'Education');
insert into authors(author_id, name) values ('b2a66341-4ef0-45bc-a00e-8f585dac788b', 'Steve Jobs');
insert into books(book_id, name, author_id, price) values  ('80f05fc5-b7c0-4717-a888-a8f881c43520','Learning Swift','b2a66341-4ef0-45bc-a00e-8f585dac788b', 50000);
insert into book_categories(book_id, category_id) values ('80f05fc5-b7c0-4717-a888-a8f881c43520', '2bd05294-e4c5-46ba-a458-ba54c799e4e3');

insert into categories(category_id, name, parent_uuid) VALUES ('2bd05294-e4c5-46ba-a458-ba54c799e4e9', 'Textbook', '2bd05294-e4c5-46ba-a458-ba54c799e4e3');
insert into categories(category_id, name, parent_uuid) VALUES ('5bd05294-e4c5-46ba-a458-ba54c799e4e4', 'Dictionary', '2bd05294-e4c5-46ba-a458-ba54c799e4e3');

/* select cat.category_id, cat.name as category_name, cat2.name as parent_category from categories as cat left join categories as cat2 on cat.parent_uuid = cat2.category_id;
