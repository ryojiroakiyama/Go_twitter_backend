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
    r_group_by_user.cnt AS "following_count",
    r_group_by_follow.cnt AS "followers_count"
  FROM account AS a
    INNER JOIN
      (SELECT
        user_id,
        COUNT(*) AS cnt
      FROM relationship
      GROUP BY user_id) AS r_group_by_user
    ON a.id = r_group_by_user.user_id
    INNER JOIN
      (SELECT
        follow_id,
        COUNT(*) AS cnt
      FROM relationship
      GROUP BY follow_id) AS r_group_by_follow
    ON a.id = r_group_by_follow.follow_id;