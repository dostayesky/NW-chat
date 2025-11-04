# Tiltfile for sa-connext (Docker Desktop)

# --- 1. Docker Builds (This part is correct) ---
docker_build(
    'connext-api-gateway:latest',
    '.', # 👈 Context คือราก
    dockerfile='infra/dockerfile/api-gateway',
    build_args={'SERVICE_NAME': 'api-gateway'}, # 👈 เพิ่ม Args
    only=['shared/', 'services/api-gateway/', 'go.mod', 'go.sum'] # 👈 เฝ้าดูไฟล์
)
docker_build(
    'connext-user-service:latest',
    '.',
    dockerfile='infra/dockerfile/user-service',
    build_args={'SERVICE_NAME': 'user-service'},
    only=['shared/', 'services/user-service/', 'go.mod', 'go.sum']
)
docker_build(
    'connext-chat-service:latest',
    '.',
    dockerfile='infra/dockerfile/chat-service',
    build_args={'SERVICE_NAME': 'chat-service'},
    only=['shared/', 'services/chat-service/', 'go.mod', 'go.sum']
)
docker_build(
    'connext-event-service:latest',
    '.',
    dockerfile='infra/dockerfile/event-service',
    build_args={'SERVICE_NAME': 'event-service'},
    only=['shared/', 'services/event-service/', 'go.mod', 'go.sum']
)
docker_build(
    'connext-notification-service:latest',
    '.',
    dockerfile='infra/dockerfile/notification-service',
    build_args={'SERVICE_NAME': 'notification-service'},
    only=['shared/', 'services/notification-service/', 'go.mod', 'go.sum']
)

# --- 2. Apply ALL k8s YAMLs ---
# (This part is correct)
k8s_yaml([
    'infra/k8s/app-config.yaml',
    'infra/k8s/app-secrets.yaml',
    'infra/k8s/db-secrets.yaml',
    'infra/k8s/postgres-db.yaml',
    'infra/k8s/pgadmin-db.yaml',

    'infra/k8s/api-gateway-service-deployment.yaml',

    'infra/k8s/user-service-deployment.yaml',

    'infra/k8s/chat-service-deployment.yaml',

    'infra/k8s/event-service-deployment.yaml',

    'infra/k8s/notification-service-deployment.yaml',
])

# --- 3. Register Resources (with Dependencies) ---

k8s_resource(
    objects=['mongodb-secret','pgadmin-secret','postgres-secret'],
    new_name='Databases Setup',
    labels='Infra'
)

k8s_resource(
    objects=['postgres-pvc'],
    new_name='Postgres Volume',
    labels='Database'
)

k8s_resource(
    objects=['app-secret','app-config'],
    new_name='Services Setup',
    labels='Infra'
)


# "ลงทะเบียน" Postgres เพื่อให้คนอื่นรอได้
k8s_resource(
    workload='postgres-deployment',
    new_name='Postgres',
    labels='Database'
)

k8s_resource(
    workload='pgadmin-deployment',
    new_name='PgAdmin',
    port_forwards=5051,
    labels='Database'
)

# API Gateway (ไม่รอ DB)
k8s_resource(workload='api-gateway-deployment',new_name='API Gateway', port_forwards=8080,labels='Gateway')

# User Service (รอ DB)
k8s_resource(
    workload='user-service-deployment', # 👈 [FIX]
    new_name='User Service',
    resource_deps=['Postgres'],
    labels='Service',
)

# Chat Service (รอ DB)
k8s_resource(
    workload='chat-service-deployment', # 👈 [FIX]
    new_name='Chat Service',
    resource_deps=['Postgres'],
    labels='Service',
)

# Event Service (รอ DB)
k8s_resource(
    workload='event-service-deployment', # 👈 [FIX]
    new_name='Event Service',
    resource_deps=['Postgres'],
    labels='Service',
)

# Notification Service (รอ DB)
k8s_resource(
    workload='notification-service-deployment', # 👈 [FIX]
    new_name='Notification Service',
    labels='Service',
)

