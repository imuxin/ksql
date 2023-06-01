SELECT
    "{ .metadata.name }" AS NAME,
    "{ .metadata.namespace }" AS NAMESPACE,
    "{ .spec.clusterIP }" AS CLUSTERIP
FROM service