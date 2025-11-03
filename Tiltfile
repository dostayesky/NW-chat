# --- 1. ตั้งค่า Kind ---
# (สำคัญ!) บอก Tilt ให้ใช้ `kind load` อัตโนมัติ
# (เปลี่ยน 'dev' ถ้าคุณใช้ชื่อ cluster อื่น)
load('ext://kind', cluster_name='dev')


# --- 2. Infrastructure (Base) ---
# YAMLs เหล่านี้จะถูก apply แค่ครั้งเดียว
# เราแยกเป็นกลุ่มๆ เพื่อให้ Tilt จัดการ dependency ได้
k8s_resource(
    'infra:secrets',
    yaml_files=[
        'k8s/postgres-secret.yaml',
        'k8s/pgadmin-secret.yaml',
        'k8s/app-secret.yaml'
    ]
)
k8s_resource(
    'infra:configs',
    yaml_file='k8s/app-config.yaml'
)
k8s_resource(
    'infra:database',
    yaml_files=[
        'k8s/postgres-db.yaml',
        'k8s/pgadmin-deployment.yaml', # (ไฟล์ pgadmin-db.yaml ของคุณ)
    ],
    resource_deps=['infra:secrets'] # 👈 รอ Secrets พร้อมก่อน
)


# --- 3. API Gateway Service ---
docker_build(
    'connext-api-gateway:latest',         # 1. ชื่อ Image
    '.',                                  # 2. (สำคัญ!) Context คือราก (Root)
    dockerfile='services/api-gateway/Dockerfile', # 3. Path ไปยัง Dockerfile
    build_args={'SERVICE_NAME': 'api-gateway'},   # 4. Build Arg ที่เราใช้
    # 5. (สำคัญ!) เฝ้าดูเฉพาะไฟล์เหล่านี้
    only=['shared/', 'services/api-gateway/', 'go.mod', 'go.sum']
)
k8s_resource(
    'api-gateway',
    yaml_files=[
        'k8s/api-gateway-deployment.yaml',
        'k8s/api-gateway-service.yaml'
    ],
    # 6. รอให้ Configs พร้อมก่อน
    resource_deps=['infra:configs']
)


# --- 4. User Service ---
docker_build(
    'connext-user-service:latest',
    '.',
    dockerfile='services/user-service/Dockerfile',
    build_args={'SERVICE_NAME': 'user-service'},
    only=['shared/', 'services/user-service/', 'go.mod', 'go.sum']
)
k8s_resource(
    'user-service',
    yaml_files=[
        'k8s/user-service-deployment.yaml',
        'k8s/user-service-service.yaml' # 👈 (ไฟล์ 'user-service.yaml' ของคุณ)
    ],
    # 7. (สำคัญ!) รอให้ DB และ Configs พร้อมก่อน
    resource_deps=['infra:database', 'infra:configs']
)


# --- 5. Chat Service ---
docker_build(
    'connext-chat-service:latest',
    '.',
    dockerfile='services/chat-service/Dockerfile',
    build_args={'SERVICE_NAME': 'chat-service'},
    only=['shared/', 'services/chat-service/', 'go.mod', 'go.sum']
)
k8s_resource(
    'chat-service',
    yaml_files=[
        'k8s/chat-service-deployment.yaml',
        'k8s/chat-service-service.yaml' # 👈 (ไฟล์ 'chat-service.yaml' ของคุณ)
    ],
    resource_deps=['infra:configs']
)


# --- 6. Event Service ---
docker_build(
    'connext-event-service:latest',
    '.',
    dockerfile='services/event-service/Dockerfile',
    build_args={'SERVICE_NAME': 'event-service'},
    only=['shared/', 'services/event-service/', 'go.mod', 'go.sum']
)
k8s_resource(
    'event-service',
    yaml_files=[
        'k8s/event-service-deployment.yaml',
        'k8s/event-service-service.yaml' # 👈 (ไฟล์ 'event-service.yaml' ของคุณ)
    ],
    resource_deps=['infra:configs']
)

# (เพิ่ม notification-service และอื่นๆ ได้ในรูปแบบเดียวกัน)

# 
# 8. ไม่ต้องใช้ k8s_resource(..., port_forwards=...)
#    เพราะเราใช้ NodePort (จาก kind-config.yaml) อยู่แล้ว
#