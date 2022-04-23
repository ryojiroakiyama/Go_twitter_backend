CREATE VIEW meta_account
AS
  SELECT
    a.username,
    a.display_name,
    a.avatar,
    a.header,
    a.note,
    a.create_at,
    r_user.cnt_user_id_rows AS "following_count",
    r_follow.cnt_follow_id_rows AS "followers_count"
  FROM account AS a
    INNER JOIN
      (SELECT
        user_id,
        COUNT(*) AS cnt_user_id_rows
      FROM relationship
      GROUP BY user_id) AS r_user
    ON a.id = r_user.user_id
    INNER JOIN
      (SELECT
        follow_id,
        COUNT(*) AS cnt_follow_id_rows
      FROM relationship
      GROUP BY follow_id) AS r_follow
    ON a.id = r_follow.follow_id;