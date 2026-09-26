SELECT
    r.oid, r.rolname, r.rolcanlogin, r.rolsuper, d.description
FROM pg_catalog.pg_roles r
    LEFT JOIN pg_catalog.pg_shdescription d
    ON d.objoid = r.oid AND d.classoid = 'pg_catalog.pg_authid'::regclass
WHERE
{% if rid %}
r.oid = {{ rid|qtLiteral(conn) }}::oid
{% else %}
(r.rolname = current_user OR (SELECT rolsuper FROM pg_catalog.pg_roles WHERE rolname = current_user))
{% endif %}
ORDER BY r.rolcanlogin, r.rolname
