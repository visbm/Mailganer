CREATE TABLE IF NOT EXISTS subscribers (
  subscriber_id serial PRIMARY KEY,
  sub_address text NOT NULL, 
  sub_name text NOT NULL,
  sub_surname text NOT NULL, 
  favourite_category text NOT NULL
);

INSERT INTO subscribers (
  sub_address, sub_name, sub_surname,
  favourite_category
) 
VALUES 
  ('visbm@mail.ru', 'Tom', 'Riddle', 'Snakes'), 
  ('ivan@mail.ru', 'Ivan', 'Ivanov', 'Cars');


CREATE TABLE IF NOT EXISTS templates (
  template_id serial PRIMARY KEY, 
  template_path text NOT NULL
);

INSERT INTO templates (template_id, template_path) 
VALUES 
  (1, 'templates/mail/hello.html'), 
  (2, 'templates/mail/hello.html');