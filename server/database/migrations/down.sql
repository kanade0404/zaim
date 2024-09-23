DROP INDEX idx_parent_genre_modified_event_child_genre_id;
DROP INDEX idx_parent_genre_modified_event_parent_genre_id;
DROP INDEX idx_parent_genre_modified_event_modified;
DROP INDEX idx_genre_modified_event_genre_id;
DROP INDEX idx_genre_modified_event_name;
DROP INDEX idx_genre_modified_event_modified;
DROP INDEX idx_genre_category_id;
DROP INDEX idx_inactive_category_category_id_inactivated_at;
DROP INDEX idx_active_category_category_id_activated;
DROP INDEX idx_category_modified_event_category_id;
DROP INDEX idx_category_modified_event_name;
DROP INDEX idx_category_modified_event_mode_id;
DROP INDEX idx_category_modified_event_modified;
DROP INDEX idx_category_category_id_user_id;
DROP INDEX idx_account_modified_event_account_id;
DROP INDEX idx_account_modified_event_name;
DROP INDEX idx_zaim_oauth_zaim_app_id;
DROP INDEX idx_zaim_app_user_id;
DROP INDEX idx_user_name;
DROP TABLE "b43_genre";
DROP TABLE "parent_genre_modified_event";
DROP TABLE "inactive_genre";
DROP TABLE "active_genre";
DROP TABLE "genre_modified_event";
DROP TABLE "genre";
DROP TABLE "inactive_category";
DROP TABLE "active_category";
DROP TABLE "category_modified_event";
DROP TABLE "category";
DROP TABLE "modes";
DROP TABLE "b43_account";
DROP TABLE "inactive_account";
DROP TABLE "active_account";
DROP TABLE "account_modified_event";
DROP TABLE "account";
DROP TABLE "enable_zaim_oauth_event";
DROP TABLE "enable_zaim_app_event";
DROP TABLE "zaim_oauth";
DROP TABLE "zaim_app";
DROP TABLE "user";