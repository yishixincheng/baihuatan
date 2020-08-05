create database bht_oauth;

CREATE TABLE IF NOT EXISTS `client_details` (
    `client_id`   VARCHAR(128) NOT NULL DEFAULT '',
    `client_secret` VARCHAR(128) NOT NULL DEFAULT '',
    `access_token_validity_seconds` INT(10) NOT NULL DEFAULT 0,
    `refresh_token_validity_seconds` INT(10) NOT NULL DEFAULT 0,
    `registered_redirect_uri` VARCHAR(128) NOT NULL DEFAULT '',
    `authorized_grant_types` VARCHAR(128) NOT NULL DEFAULT '',
    PRIMARY KEY (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into `client_details` (`client_id`, `client_secret`, `access_token_validity_seconds`, `refresh_token_validity_seconds`, `registered_redirect_uri`, `authorized_grant_types`) value (`clientId`, 'clientSecret', 1800, 18000, 'http://localhost', '["password","refresh_token"]');

