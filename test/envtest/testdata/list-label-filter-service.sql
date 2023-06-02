SELECT * FROM service
    LABEL app = nginx
    LABEL app == nginx
    LABEL "xxx" not exists
    LABEL "app" exists
    LABEL "version" in (0, 1, 2)
    LABEL "version" >= 0
    LABEL "version" <= 100
    LABEL "version" != 10
