CREATE VIEW meta_account
AS
  SELECT
    a.id,
    a.username,
    a.password_hash,
    a.display_name,
    a.avatar,
    a.header,
    a.note,
    a.create_at,
    CASE WHEN r_group_by_user.cnt IS NULL
         THEN 0
         ELSE r_group_by_user.cnt
    END AS "following_count",
    CASE WHEN r_group_by_follow.cnt IS NULL
         THEN 0
         ELSE r_group_by_follow.cnt
    END AS "followers_count"
  FROM account AS a
    LEFT OUTER JOIN
      (SELECT
        user_id,
        COUNT(*) AS cnt
      FROM relationship
      GROUP BY user_id) AS r_group_by_user
    ON a.id = r_group_by_user.user_id
    LEFT OUTER JOIN
      (SELECT
        follow_id,
        COUNT(*) AS cnt
      FROM relationship
      GROUP BY follow_id) AS r_group_by_follow
    ON a.id = r_group_by_follow.follow_id;
