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


select
    b.book_id as book_id,
    b.name as book_name,
    b.price as book_price,
    b.created_at as b_created_at,
    b.updated_at as b_updated_at,
    c.category_id as category_id,
    c.name as category_name,
    c.parent_uuid as parent_uuid,
    c.created_at c_created_at,
    c.updated_at c_updated_at,
    a.author_id as author_id,
    a.name as author_name,
    a.created_at as a_created_at,
    a.updated_at as a_updated_at
from
    book_categories
join books b on book_categories.book_id = b.book_id
join categories c on book_categories.category_id = c.category_id
join authors a on b.author_id = a.author_id
where b.deleted_at is null
;

select
    c.category_id as category_id,
    c.name as category_name,
    c.parent_uuid as parent_uuid,
    c.created_at c_created_at,
    c.updated_at c_updated_at
from
    book_categories
        join books b on book_categories.book_id = b.book_id
        join categories c on book_categories.category_id = c.category_id
where b.deleted_at is null and b.book_id = '3cbecabd-b154-442c-84a8-1e587048542d'
;

select b.name, c.name
from book_categories
join books b on book_categories.book_id = b.book_id
join categories c on book_categories.category_id = c.category_id
where b.book_id = '80f05fc5-b7c0-4717-a888-a8f881c43520'
;

/* select cat.category_id, cat.name as category_name, cat2.name as parent_category from categories as cat left join categories as cat2 on cat.parent_uuid = cat2.category_id;