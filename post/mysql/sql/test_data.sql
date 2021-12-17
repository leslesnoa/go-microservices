DROP TABLE IF EXISTS `posts`;

create table IF not exists `posts`
(
 `id`               INT(20) NOT NULL AUTO_INCREMENT,
 `title`            VARCHAR(50) NOT NULL,
 `text`             VARCHAR(50) NOT NULL,
    PRIMARY KEY (`id`)
);

INSERT INTO posts (title, text) VALUES
    ('title1', 'text1'),
    ('title2', 'text2'),
    ('title3', 'text3');