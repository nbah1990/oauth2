/*!40101 SET @OLD_CHARACTER_SET_CLIENT = @@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;
/*!40101 SET @OLD_SQL_MODE = @@SQL_MODE, SQL_MODE = 'NO_AUTO_VALUE_ON_ZERO' */;

CREATE TABLE IF NOT EXISTS `oauth_access_tokens`
(
    `id`         varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `user_id`    VARCHAR(36)                                      DEFAULT '',
    `client_id`  VARCHAR(36)                             NOT NULL,
    `name`       varchar(255) COLLATE utf8mb4_unicode_ci          DEFAULT NULL,
    `scopes`     text COLLATE utf8mb4_unicode_ci,
    `revoked`    tinyint(1)                              NOT NULL DEFAULT 0,
    `created_at` timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `expires_at` datetime                                NOT NULL,
    PRIMARY KEY (`id`),
    KEY `oauth_access_tokens_user_id_index` (`user_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `oauth_clients`
(
    `id`                     VARCHAR(36)                             NOT NULL,
    `user_id`                VARCHAR(36)                                      DEFAULT '',
    `name`                   varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `secret`                 varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `scopes`                 text COLLATE utf8mb4_unicode_ci,
    `redirect`               text COLLATE utf8mb4_unicode_ci         NULL DEFAULT NULL,
    `personal_access_client` tinyint(1)                              NOT NULL DEFAULT 0,
    `password_client`        tinyint(1)                              NOT NULL DEFAULT 1,
    `revoked`                tinyint(1)                              NOT NULL DEFAULT 0,
    `created_at`             timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`             timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `oauth_clients_user_id_index` (`user_id`),
    UNIQUE KEY `name_unique` (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `oauth_refresh_tokens`
(
    `id`              varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `access_token_id` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    `revoked`         tinyint(1)                              NOT NULL,
    `expires_at`      datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `oauth_refresh_tokens_access_token_id_index` (`access_token_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

/*!40101 SET SQL_MODE = IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS = IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT = @OLD_CHARACTER_SET_CLIENT */;