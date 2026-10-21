-- +goose Up
-- Same schema the Laravel migrations produced, so a database that already
-- has these tables keeps working. IF NOT EXISTS makes this a no-op there.
CREATE TABLE IF NOT EXISTS `irma_meet_sessions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `irma_session_id` char(48) NOT NULL,
  `meeting_name` varchar(255) NOT NULL,
  `hoster_name` varchar(255) NOT NULL,
  `hoster_email_address` varchar(255) NOT NULL,
  `start_time` datetime NOT NULL,
  `invitation_note` varchar(255) DEFAULT '',
  `bbb_session_id` varchar(255) NOT NULL,
  `bbb_moderator_password` varchar(255) DEFAULT NULL,
  `bbb_attendee_password` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `meeting_type` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `irma_meet_sessions_irma_session_id_unique` (`irma_session_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `irma_meet_participants` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `irma_session_id` char(48) NOT NULL,
  `email_address` varchar(255) NOT NULL,
  `authentication` int NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS `irma_meet_participants`;
DROP TABLE IF EXISTS `irma_meet_sessions`;
