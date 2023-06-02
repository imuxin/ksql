SELECT * FROM service
    WHERE
        "{ .metadata.namespace }" = "ksql" OR
        "{ .metadata.namespace }" = "ksql-system" AND
        "{ .spec.clusterIP }" in ("10.0.0.136", "10.0.0.137") AND
        "{ .spec.clusterIP }" not in ("10.0.0.1", "10.0.0.2") AND
        "{ .spec.ports[0].port }" != 443 AND
        "{ .spec.ports[0].port }" > 22 AND
        "{ .spec.ports[0].port }" < 8080
