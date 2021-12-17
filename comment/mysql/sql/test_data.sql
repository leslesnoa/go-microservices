DROP TABLE IF EXISTS `comments`;

create table IF not exists `comments`
(
 `id`               INT(20) NOT NULL AUTO_INCREMENT,
 `comment_by_post_id`  INT(20) NOT NULL,
 `content`             VARCHAR(50) NOT NULL,
    PRIMARY KEY (`id`)
);

INSERT INTO comments (comment_by_post_id, content) VALUES
    (1, 'comment1'),
    (2, 'comment2'),
    (3, 'comment3');