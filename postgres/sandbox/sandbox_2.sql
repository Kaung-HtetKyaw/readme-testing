BEGIN;

 
INSERT INTO sandbox (id, name) VALUES ('3d6868ed-4a22-4c6b-a740-359d5fc2816d', 'sandbox_2');


INSERT INTO sandbox_cluster_info (sandbox_id, name, version, platform) VALUES ('3d6868ed-4a22-4c6b-a740-359d5fc2816d', 'sandbox_2', '1.31.9-eks-5d4a308', 'linux/amd64');


INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'cd2942c1-563b-4738-ac8b-39f94494d142',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_1',
'9104376',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "cd2942c1-563b-4738-ac8b-39f94494d142", "name": "object_1", "resourceVersion": "9104376", "labels": {"kubernetes.io/metadata.name": "object_1"}, "managedFields": [], "creationTimestamp": "2025-05-19T18:01:05.454672Z"}}',
'2025-05-19T18:01:05.454672Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'c7ebe891-06ad-429d-bdb1-ebf2f9b88f90',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_2',
'6969324',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "c7ebe891-06ad-429d-bdb1-ebf2f9b88f90", "name": "object_2", "resourceVersion": "6969324", "labels": {"kubernetes.io/metadata.name": "object_2"}, "managedFields": [], "creationTimestamp": "2025-06-27T18:01:05.454777Z"}}',
'2025-06-27T18:01:05.454777Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'9b637b4d-a55e-4db4-aa45-3e7a30ccb6f3',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_3',
'4726335',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "9b637b4d-a55e-4db4-aa45-3e7a30ccb6f3", "name": "object_3", "resourceVersion": "4726335", "labels": {"kubernetes.io/metadata.name": "object_3"}, "managedFields": [], "creationTimestamp": "2025-05-09T18:01:05.454911Z"}}',
'2025-05-09T18:01:05.454911Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'7fd440dd-802e-4939-bd74-e4ebe3fd4c98',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_4',
'2191421',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "7fd440dd-802e-4939-bd74-e4ebe3fd4c98", "name": "object_4", "resourceVersion": "2191421", "labels": {"kubernetes.io/metadata.name": "object_4"}, "managedFields": [], "creationTimestamp": "2025-07-02T18:01:05.454956Z"}}',
'2025-07-02T18:01:05.454956Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'337fa1f3-b0a2-401a-8551-e32c0ae4df17',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_5',
'2989822',
'ClusterRole',
'{"kind": "ClusterRole", "metadata": {"uid": "337fa1f3-b0a2-401a-8551-e32c0ae4df17", "name": "object_5", "resourceVersion": "2989822", "labels": {"kubernetes.io/metadata.name": "object_5"}, "managedFields": [], "creationTimestamp": "2025-06-25T18:01:05.455016Z"}}',
'2025-06-25T18:01:05.455016Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'ba005a3c-0dfb-4af8-ab26-dcd9120ce6bf',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_6',
'7373663',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "ba005a3c-0dfb-4af8-ab26-dcd9120ce6bf", "name": "object_6", "resourceVersion": "7373663", "labels": {"kubernetes.io/metadata.name": "object_6"}, "managedFields": [], "creationTimestamp": "2025-06-12T18:01:05.455109Z"}}',
'2025-06-12T18:01:05.455109Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'ca7e2aa2-82ed-4e59-b0b1-f7cc89da70ce',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_7',
'2912935',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "ca7e2aa2-82ed-4e59-b0b1-f7cc89da70ce", "name": "object_7", "resourceVersion": "2912935", "labels": {"kubernetes.io/metadata.name": "object_7"}, "managedFields": [], "creationTimestamp": "2025-06-03T18:01:05.455372Z"}}',
'2025-06-03T18:01:05.455372Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'4a148dca-f339-4b77-8a2f-056e3c68ea99',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_8',
'2855147',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "4a148dca-f339-4b77-8a2f-056e3c68ea99", "name": "object_8", "resourceVersion": "2855147", "labels": {"kubernetes.io/metadata.name": "object_8"}, "managedFields": [], "creationTimestamp": "2025-06-12T18:01:05.455439Z"}}',
'2025-06-12T18:01:05.455439Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'a3ce1327-e7d6-4311-8a75-b87a3e0ea5d3',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_9',
'2400004',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "a3ce1327-e7d6-4311-8a75-b87a3e0ea5d3", "name": "object_9", "resourceVersion": "2400004", "labels": {"kubernetes.io/metadata.name": "object_9"}, "managedFields": [], "creationTimestamp": "2025-05-24T18:01:05.455486Z"}}',
'2025-05-24T18:01:05.455486Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'3b81a091-1217-4ed6-8c20-296e833b77a5',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_10',
'3806157',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "3b81a091-1217-4ed6-8c20-296e833b77a5", "name": "object_10", "resourceVersion": "3806157", "labels": {"kubernetes.io/metadata.name": "object_10"}, "managedFields": [], "creationTimestamp": "2025-07-01T18:01:05.455535Z"}}',
'2025-07-01T18:01:05.455535Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'5b3f22f6-65d7-4f71-84af-5af4b14e7d41',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_11',
'7481474',
'ClusterRole',
'{"kind": "ClusterRole", "metadata": {"uid": "5b3f22f6-65d7-4f71-84af-5af4b14e7d41", "name": "object_11", "resourceVersion": "7481474", "labels": {"kubernetes.io/metadata.name": "object_11"}, "managedFields": [], "creationTimestamp": "2025-06-06T18:01:05.455585Z"}}',
'2025-06-06T18:01:05.455585Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'b50aac26-fc1d-407f-8b73-134537472293',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_12',
'8966313',
'ClusterRole',
'{"kind": "ClusterRole", "metadata": {"uid": "b50aac26-fc1d-407f-8b73-134537472293", "name": "object_12", "resourceVersion": "8966313", "labels": {"kubernetes.io/metadata.name": "object_12"}, "managedFields": [], "creationTimestamp": "2025-05-19T18:01:05.455737Z"}}',
'2025-05-19T18:01:05.455737Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'a5129313-5708-43fe-a3d0-11d834441f60',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_13',
'1016324',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "a5129313-5708-43fe-a3d0-11d834441f60", "name": "object_13", "resourceVersion": "1016324", "labels": {"kubernetes.io/metadata.name": "object_13"}, "managedFields": [], "creationTimestamp": "2025-06-16T18:01:05.455810Z"}}',
'2025-06-16T18:01:05.455810Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'a8e6a435-4314-4da8-a19f-69ac92b06f07',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_14',
'2550891',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "a8e6a435-4314-4da8-a19f-69ac92b06f07", "name": "object_14", "resourceVersion": "2550891", "labels": {"kubernetes.io/metadata.name": "object_14"}, "managedFields": [], "creationTimestamp": "2025-05-23T18:01:05.455958Z"}}',
'2025-05-23T18:01:05.455958Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'7dc10f2f-a6dc-4a8e-ac06-5845c7225748',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_15',
'2701241',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "7dc10f2f-a6dc-4a8e-ac06-5845c7225748", "name": "object_15", "resourceVersion": "2701241", "labels": {"kubernetes.io/metadata.name": "object_15"}, "managedFields": [], "creationTimestamp": "2025-07-06T18:01:05.456005Z"}}',
'2025-07-06T18:01:05.456005Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'4282fcfb-f655-4179-bdfd-869a49f421eb',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_16',
'6878063',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "4282fcfb-f655-4179-bdfd-869a49f421eb", "name": "object_16", "resourceVersion": "6878063", "labels": {"kubernetes.io/metadata.name": "object_16"}, "managedFields": [], "creationTimestamp": "2025-05-19T18:01:05.456037Z"}}',
'2025-05-19T18:01:05.456037Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'2cd715cf-1b72-47a5-984e-c33ef73ef046',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_17',
'7395149',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "2cd715cf-1b72-47a5-984e-c33ef73ef046", "name": "object_17", "resourceVersion": "7395149", "labels": {"kubernetes.io/metadata.name": "object_17"}, "managedFields": [], "creationTimestamp": "2025-05-19T18:01:05.456080Z"}}',
'2025-05-19T18:01:05.456080Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'21f2df90-f440-4164-96ba-9416c8ac0b32',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_18',
'1325684',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "21f2df90-f440-4164-96ba-9416c8ac0b32", "name": "object_18", "resourceVersion": "1325684", "labels": {"kubernetes.io/metadata.name": "object_18"}, "managedFields": [], "creationTimestamp": "2025-06-28T18:01:05.456128Z"}}',
'2025-06-28T18:01:05.456128Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'22d3886e-0396-4fd2-a713-59cf776ce59c',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_19',
'9682502',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "22d3886e-0396-4fd2-a713-59cf776ce59c", "name": "object_19", "resourceVersion": "9682502", "labels": {"kubernetes.io/metadata.name": "object_19"}, "managedFields": [], "creationTimestamp": "2025-06-25T18:01:05.456184Z"}}',
'2025-06-25T18:01:05.456184Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'015cda50-499f-48a2-9009-bcf1a21270ec',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_20',
'5616317',
'ClusterRole',
'{"kind": "ClusterRole", "metadata": {"uid": "015cda50-499f-48a2-9009-bcf1a21270ec", "name": "object_20", "resourceVersion": "5616317", "labels": {"kubernetes.io/metadata.name": "object_20"}, "managedFields": [], "creationTimestamp": "2025-05-28T18:01:05.456355Z"}}',
'2025-05-28T18:01:05.456355Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'da84763d-5612-4887-b1d7-c88181cbe672',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_21',
'8243298',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "da84763d-5612-4887-b1d7-c88181cbe672", "name": "object_21", "resourceVersion": "8243298", "labels": {"kubernetes.io/metadata.name": "object_21"}, "managedFields": [], "creationTimestamp": "2025-06-24T18:01:05.456434Z"}}',
'2025-06-24T18:01:05.456434Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'69cf87d4-169b-4108-8a5c-f5af78b9b29a',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_22',
'7328953',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "69cf87d4-169b-4108-8a5c-f5af78b9b29a", "name": "object_22", "resourceVersion": "7328953", "labels": {"kubernetes.io/metadata.name": "object_22"}, "managedFields": [], "creationTimestamp": "2025-06-22T18:01:05.456485Z"}}',
'2025-06-22T18:01:05.456485Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'87e8c15a-c551-4c40-8abb-740804bcf94a',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_23',
'4948550',
'ClusterRole',
'{"kind": "ClusterRole", "metadata": {"uid": "87e8c15a-c551-4c40-8abb-740804bcf94a", "name": "object_23", "resourceVersion": "4948550", "labels": {"kubernetes.io/metadata.name": "object_23"}, "managedFields": [], "creationTimestamp": "2025-06-08T18:01:05.456525Z"}}',
'2025-06-08T18:01:05.456525Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'5e6f6e06-0949-4627-ae2c-362f143dcc3a',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_24',
'2088047',
'ClusterRole',
'{"kind": "ClusterRole", "metadata": {"uid": "5e6f6e06-0949-4627-ae2c-362f143dcc3a", "name": "object_24", "resourceVersion": "2088047", "labels": {"kubernetes.io/metadata.name": "object_24"}, "managedFields": [], "creationTimestamp": "2025-05-19T18:01:05.456567Z"}}',
'2025-05-19T18:01:05.456567Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'6cc196ee-4d5d-4e92-8a43-d84ed375b11a',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_25',
'7357383',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "6cc196ee-4d5d-4e92-8a43-d84ed375b11a", "name": "object_25", "resourceVersion": "7357383", "labels": {"kubernetes.io/metadata.name": "object_25"}, "managedFields": [], "creationTimestamp": "2025-06-04T18:01:05.456606Z"}}',
'2025-06-04T18:01:05.456606Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'5e5b523e-632b-4c0a-ab56-d18b32943d81',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_26',
'7512157',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "5e5b523e-632b-4c0a-ab56-d18b32943d81", "name": "object_26", "resourceVersion": "7512157", "labels": {"kubernetes.io/metadata.name": "object_26"}, "managedFields": [], "creationTimestamp": "2025-06-09T18:01:05.456838Z"}}',
'2025-06-09T18:01:05.456838Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'5e8a9592-3e69-4ee7-b44d-de76d524e85f',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_27',
'4821149',
'ClusterRole',
'{"kind": "ClusterRole", "metadata": {"uid": "5e8a9592-3e69-4ee7-b44d-de76d524e85f", "name": "object_27", "resourceVersion": "4821149", "labels": {"kubernetes.io/metadata.name": "object_27"}, "managedFields": [], "creationTimestamp": "2025-06-07T18:01:05.456879Z"}}',
'2025-06-07T18:01:05.456879Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'89086361-b41a-401a-8bb4-c9c8f8dac9c6',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_28',
'1932483',
'ClusterRoleBinding',
'{"kind": "ClusterRoleBinding", "metadata": {"uid": "89086361-b41a-401a-8bb4-c9c8f8dac9c6", "name": "object_28", "resourceVersion": "1932483", "labels": {"kubernetes.io/metadata.name": "object_28"}, "managedFields": [], "creationTimestamp": "2025-05-10T18:01:05.457021Z"}}',
'2025-05-10T18:01:05.457021Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'72760778-f5cd-4a65-acad-2713186ed4d6',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_29',
'6564862',
'Namespace',
'{"kind": "Namespace", "metadata": {"uid": "72760778-f5cd-4a65-acad-2713186ed4d6", "name": "object_29", "resourceVersion": "6564862", "labels": {"kubernetes.io/metadata.name": "object_29"}, "managedFields": [], "creationTimestamp": "2025-06-21T18:01:05.457108Z"}}',
'2025-06-21T18:01:05.457108Z'
);

INSERT INTO sandbox_object (id, sandbox_id, name, resource_version, kind, raw, created_at) VALUES (
'00ac4a50-25cf-433d-a17b-cabf1327e54e',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'object_30',
'6196851',
'ClusterRole',
'{"kind": "ClusterRole", "metadata": {"uid": "00ac4a50-25cf-433d-a17b-cabf1327e54e", "name": "object_30", "resourceVersion": "6196851", "labels": {"kubernetes.io/metadata.name": "object_30"}, "managedFields": [], "creationTimestamp": "2025-05-15T18:01:05.457180Z"}}',
'2025-05-15T18:01:05.457180Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'd26058d2-d87e-4bd7-b757-3a1d7ed0c70e',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'dev-pod-28',
'dev',
'8730613',
'{"kind": "Pod", "metadata": {"uid": "d26058d2-d87e-4bd7-b757-3a1d7ed0c70e", "name": "dev-pod-28", "namespace": "dev", "resourceVersion": "8730613", "labels": {"kubernetes.io/metadata.name": "dev-pod-28"}, "managedFields": [], "creationTimestamp": "2025-07-07T18:01:05.457410Z"}}',
'2025-07-07T18:01:05.457410Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'ebc29e5e-21f3-4347-b90b-80cfc97e267d',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'kube-system-configmap-30',
'kube-system',
'4225193',
'{"kind": "ConfigMap", "metadata": {"uid": "ebc29e5e-21f3-4347-b90b-80cfc97e267d", "name": "kube-system-configmap-30", "namespace": "kube-system", "resourceVersion": "4225193", "labels": {"kubernetes.io/metadata.name": "kube-system-configmap-30"}, "managedFields": [], "creationTimestamp": "2025-05-13T18:01:05.457476Z"}}',
'2025-05-13T18:01:05.457476Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'54513302-2acb-4dd0-8d6c-4d6c83544ff3',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'monitoring-service-70',
'monitoring',
'9903642',
'{"kind": "Service", "metadata": {"uid": "54513302-2acb-4dd0-8d6c-4d6c83544ff3", "name": "monitoring-service-70", "namespace": "monitoring", "resourceVersion": "9903642", "labels": {"kubernetes.io/metadata.name": "monitoring-service-70"}, "managedFields": [], "creationTimestamp": "2025-06-27T18:01:05.457533Z"}}',
'2025-06-27T18:01:05.457533Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'aa01888b-5d5c-46fc-a5ff-90d6256290ee',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'kube-system-configmap-57',
'kube-system',
'4092175',
'{"kind": "ConfigMap", "metadata": {"uid": "aa01888b-5d5c-46fc-a5ff-90d6256290ee", "name": "kube-system-configmap-57", "namespace": "kube-system", "resourceVersion": "4092175", "labels": {"kubernetes.io/metadata.name": "kube-system-configmap-57"}, "managedFields": [], "creationTimestamp": "2025-07-02T18:01:05.457587Z"}}',
'2025-07-02T18:01:05.457587Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'43496452-5869-4ab6-bc04-c718aa677080',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'kube-system-service-62',
'kube-system',
'8842252',
'{"kind": "Service", "metadata": {"uid": "43496452-5869-4ab6-bc04-c718aa677080", "name": "kube-system-service-62", "namespace": "kube-system", "resourceVersion": "8842252", "labels": {"kubernetes.io/metadata.name": "kube-system-service-62"}, "managedFields": [], "creationTimestamp": "2025-06-05T18:01:05.457677Z"}}',
'2025-06-05T18:01:05.457677Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'9967c542-a22a-4306-a52c-edc154b79996',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'team-b-pod-93',
'team-b',
'7871317',
'{"kind": "Pod", "metadata": {"uid": "9967c542-a22a-4306-a52c-edc154b79996", "name": "team-b-pod-93", "namespace": "team-b", "resourceVersion": "7871317", "labels": {"kubernetes.io/metadata.name": "team-b-pod-93"}, "managedFields": [], "creationTimestamp": "2025-06-07T18:01:05.457740Z"}}',
'2025-06-07T18:01:05.457740Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'ffda4e3f-b954-4637-9226-be5aadc38879',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'staging-service-74',
'staging',
'1379102',
'{"kind": "Service", "metadata": {"uid": "ffda4e3f-b954-4637-9226-be5aadc38879", "name": "staging-service-74", "namespace": "staging", "resourceVersion": "1379102", "labels": {"kubernetes.io/metadata.name": "staging-service-74"}, "managedFields": [], "creationTimestamp": "2025-06-26T18:01:05.457821Z"}}',
'2025-06-26T18:01:05.457821Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'237cc8f6-09bc-4adf-80b6-1336d322e087',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'kube-system-configmap-2',
'kube-system',
'6747181',
'{"kind": "ConfigMap", "metadata": {"uid": "237cc8f6-09bc-4adf-80b6-1336d322e087", "name": "kube-system-configmap-2", "namespace": "kube-system", "resourceVersion": "6747181", "labels": {"kubernetes.io/metadata.name": "kube-system-configmap-2"}, "managedFields": [], "creationTimestamp": "2025-06-26T18:01:05.457886Z"}}',
'2025-06-26T18:01:05.457886Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'c69b3dd1-6883-4714-aee0-33fb3269c53d',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'monitoring-service-21',
'monitoring',
'3305419',
'{"kind": "Service", "metadata": {"uid": "c69b3dd1-6883-4714-aee0-33fb3269c53d", "name": "monitoring-service-21", "namespace": "monitoring", "resourceVersion": "3305419", "labels": {"kubernetes.io/metadata.name": "monitoring-service-21"}, "managedFields": [], "creationTimestamp": "2025-06-24T18:01:05.457942Z"}}',
'2025-06-24T18:01:05.457942Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'e9ece1d6-8240-4708-93bb-9b4cd5456f3f',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'prometheus-configmap-98',
'prometheus',
'4547218',
'{"kind": "ConfigMap", "metadata": {"uid": "e9ece1d6-8240-4708-93bb-9b4cd5456f3f", "name": "prometheus-configmap-98", "namespace": "prometheus", "resourceVersion": "4547218", "labels": {"kubernetes.io/metadata.name": "prometheus-configmap-98"}, "managedFields": [], "creationTimestamp": "2025-06-07T18:01:05.457995Z"}}',
'2025-06-07T18:01:05.457995Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'e677a0cc-a778-4dfb-a5e5-c2533ac2bbba',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'monitoring-deployment-97',
'monitoring',
'8253560',
'{"kind": "Deployment", "metadata": {"uid": "e677a0cc-a778-4dfb-a5e5-c2533ac2bbba", "name": "monitoring-deployment-97", "namespace": "monitoring", "resourceVersion": "8253560", "labels": {"kubernetes.io/metadata.name": "monitoring-deployment-97"}, "managedFields": [], "creationTimestamp": "2025-05-15T18:01:05.458062Z"}}',
'2025-05-15T18:01:05.458062Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'84775c8f-7bfe-402c-8ba4-fe562c607d6d',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'monitoring-deployment-53',
'monitoring',
'1345350',
'{"kind": "Deployment", "metadata": {"uid": "84775c8f-7bfe-402c-8ba4-fe562c607d6d", "name": "monitoring-deployment-53", "namespace": "monitoring", "resourceVersion": "1345350", "labels": {"kubernetes.io/metadata.name": "monitoring-deployment-53"}, "managedFields": [], "creationTimestamp": "2025-06-03T18:01:05.458113Z"}}',
'2025-06-03T18:01:05.458113Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'5fc076a6-c1e7-42e9-88a4-ca4f8cd3ec04',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'grafana-deployment-98',
'grafana',
'6662454',
'{"kind": "Deployment", "metadata": {"uid": "5fc076a6-c1e7-42e9-88a4-ca4f8cd3ec04", "name": "grafana-deployment-98", "namespace": "grafana", "resourceVersion": "6662454", "labels": {"kubernetes.io/metadata.name": "grafana-deployment-98"}, "managedFields": [], "creationTimestamp": "2025-06-20T18:01:05.458185Z"}}',
'2025-06-20T18:01:05.458185Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'e8aa64a4-62ca-450d-ab5d-df854705b273',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'team-a-configmap-44',
'team-a',
'4615817',
'{"kind": "ConfigMap", "metadata": {"uid": "e8aa64a4-62ca-450d-ab5d-df854705b273", "name": "team-a-configmap-44", "namespace": "team-a", "resourceVersion": "4615817", "labels": {"kubernetes.io/metadata.name": "team-a-configmap-44"}, "managedFields": [], "creationTimestamp": "2025-05-29T18:01:05.458389Z"}}',
'2025-05-29T18:01:05.458389Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'4470a2eb-a70f-4f17-85b2-0eba07f865d4',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'logging-deployment-94',
'logging',
'8443933',
'{"kind": "Deployment", "metadata": {"uid": "4470a2eb-a70f-4f17-85b2-0eba07f865d4", "name": "logging-deployment-94", "namespace": "logging", "resourceVersion": "8443933", "labels": {"kubernetes.io/metadata.name": "logging-deployment-94"}, "managedFields": [], "creationTimestamp": "2025-06-08T18:01:05.458470Z"}}',
'2025-06-08T18:01:05.458470Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'46426f10-594c-4f96-b5b4-06b45dd8e26f',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'prometheus-pod-19',
'prometheus',
'5154339',
'{"kind": "Pod", "metadata": {"uid": "46426f10-594c-4f96-b5b4-06b45dd8e26f", "name": "prometheus-pod-19", "namespace": "prometheus", "resourceVersion": "5154339", "labels": {"kubernetes.io/metadata.name": "prometheus-pod-19"}, "managedFields": [], "creationTimestamp": "2025-06-08T18:01:05.458527Z"}}',
'2025-06-08T18:01:05.458527Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'ad1316bf-c6f8-4f7c-a2dd-c3227ae59f18',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'team-a-service-24',
'team-a',
'4416766',
'{"kind": "Service", "metadata": {"uid": "ad1316bf-c6f8-4f7c-a2dd-c3227ae59f18", "name": "team-a-service-24", "namespace": "team-a", "resourceVersion": "4416766", "labels": {"kubernetes.io/metadata.name": "team-a-service-24"}, "managedFields": [], "creationTimestamp": "2025-06-03T18:01:05.458579Z"}}',
'2025-06-03T18:01:05.458579Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'98f571e1-710f-413b-8a3a-2ac034a326e2',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'team-a-service-54',
'team-a',
'8444753',
'{"kind": "Service", "metadata": {"uid": "98f571e1-710f-413b-8a3a-2ac034a326e2", "name": "team-a-service-54", "namespace": "team-a", "resourceVersion": "8444753", "labels": {"kubernetes.io/metadata.name": "team-a-service-54"}, "managedFields": [], "creationTimestamp": "2025-06-23T18:01:05.458652Z"}}',
'2025-06-23T18:01:05.458652Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'3443733b-35f0-4483-8793-baacc77d8a83',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'monitoring-deployment-51',
'monitoring',
'2991719',
'{"kind": "Deployment", "metadata": {"uid": "3443733b-35f0-4483-8793-baacc77d8a83", "name": "monitoring-deployment-51", "namespace": "monitoring", "resourceVersion": "2991719", "labels": {"kubernetes.io/metadata.name": "monitoring-deployment-51"}, "managedFields": [], "creationTimestamp": "2025-05-15T18:01:05.458708Z"}}',
'2025-05-15T18:01:05.458708Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'3d6205c0-6c8c-4188-932e-bde21796d653',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'team-b-configmap-12',
'team-b',
'8528619',
'{"kind": "ConfigMap", "metadata": {"uid": "3d6205c0-6c8c-4188-932e-bde21796d653", "name": "team-b-configmap-12", "namespace": "team-b", "resourceVersion": "8528619", "labels": {"kubernetes.io/metadata.name": "team-b-configmap-12"}, "managedFields": [], "creationTimestamp": "2025-05-25T18:01:05.458773Z"}}',
'2025-05-25T18:01:05.458773Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'01e3e36a-8d04-4ecb-b1fe-a37ad25bcb4e',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'logging-service-27',
'logging',
'1018710',
'{"kind": "Service", "metadata": {"uid": "01e3e36a-8d04-4ecb-b1fe-a37ad25bcb4e", "name": "logging-service-27", "namespace": "logging", "resourceVersion": "1018710", "labels": {"kubernetes.io/metadata.name": "logging-service-27"}, "managedFields": [], "creationTimestamp": "2025-06-09T18:01:05.458827Z"}}',
'2025-06-09T18:01:05.458827Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'7351c1a8-6976-4940-b722-e53586288c59',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'ingress-nginx-deployment-60',
'ingress-nginx',
'9595753',
'{"kind": "Deployment", "metadata": {"uid": "7351c1a8-6976-4940-b722-e53586288c59", "name": "ingress-nginx-deployment-60", "namespace": "ingress-nginx", "resourceVersion": "9595753", "labels": {"kubernetes.io/metadata.name": "ingress-nginx-deployment-60"}, "managedFields": [], "creationTimestamp": "2025-05-12T18:01:05.458959Z"}}',
'2025-05-12T18:01:05.458959Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'3b94b68a-609c-4903-a6c9-b3ca8989a1f0',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'logging-configmap-86',
'logging',
'2029969',
'{"kind": "ConfigMap", "metadata": {"uid": "3b94b68a-609c-4903-a6c9-b3ca8989a1f0", "name": "logging-configmap-86", "namespace": "logging", "resourceVersion": "2029969", "labels": {"kubernetes.io/metadata.name": "logging-configmap-86"}, "managedFields": [], "creationTimestamp": "2025-06-25T18:01:05.459034Z"}}',
'2025-06-25T18:01:05.459034Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'c7ed45bd-f0c8-46d3-9c04-6baec2268943',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'monitoring-deployment-30',
'monitoring',
'5542501',
'{"kind": "Deployment", "metadata": {"uid": "c7ed45bd-f0c8-46d3-9c04-6baec2268943", "name": "monitoring-deployment-30", "namespace": "monitoring", "resourceVersion": "5542501", "labels": {"kubernetes.io/metadata.name": "monitoring-deployment-30"}, "managedFields": [], "creationTimestamp": "2025-05-27T18:01:05.459083Z"}}',
'2025-05-27T18:01:05.459083Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'0903f3a9-bacc-41eb-a321-fd23749d0d94',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'logging-configmap-46',
'logging',
'1406916',
'{"kind": "ConfigMap", "metadata": {"uid": "0903f3a9-bacc-41eb-a321-fd23749d0d94", "name": "logging-configmap-46", "namespace": "logging", "resourceVersion": "1406916", "labels": {"kubernetes.io/metadata.name": "logging-configmap-46"}, "managedFields": [], "creationTimestamp": "2025-06-25T18:01:05.459145Z"}}',
'2025-06-25T18:01:05.459145Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'c47ddf53-616a-4345-a09e-8647e1c54c0b',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'prometheus-service-31',
'prometheus',
'3289825',
'{"kind": "Service", "metadata": {"uid": "c47ddf53-616a-4345-a09e-8647e1c54c0b", "name": "prometheus-service-31", "namespace": "prometheus", "resourceVersion": "3289825", "labels": {"kubernetes.io/metadata.name": "prometheus-service-31"}, "managedFields": [], "creationTimestamp": "2025-05-09T18:01:05.459194Z"}}',
'2025-05-09T18:01:05.459194Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'a2ce1d48-a933-4411-8e2e-acb549b23297',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'staging-deployment-21',
'staging',
'8280468',
'{"kind": "Deployment", "metadata": {"uid": "a2ce1d48-a933-4411-8e2e-acb549b23297", "name": "staging-deployment-21", "namespace": "staging", "resourceVersion": "8280468", "labels": {"kubernetes.io/metadata.name": "staging-deployment-21"}, "managedFields": [], "creationTimestamp": "2025-06-10T18:01:05.459319Z"}}',
'2025-06-10T18:01:05.459319Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'45d9459f-3a87-4d82-a8e5-82f7496bafb9',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'monitoring-deployment-57',
'monitoring',
'8211819',
'{"kind": "Deployment", "metadata": {"uid": "45d9459f-3a87-4d82-a8e5-82f7496bafb9", "name": "monitoring-deployment-57", "namespace": "monitoring", "resourceVersion": "8211819", "labels": {"kubernetes.io/metadata.name": "monitoring-deployment-57"}, "managedFields": [], "creationTimestamp": "2025-06-26T18:01:05.459417Z"}}',
'2025-06-26T18:01:05.459417Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'29334e4e-aca9-44a3-8caf-f0ab649bac6d',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'argocd-configmap-81',
'argocd',
'2817286',
'{"kind": "ConfigMap", "metadata": {"uid": "29334e4e-aca9-44a3-8caf-f0ab649bac6d", "name": "argocd-configmap-81", "namespace": "argocd", "resourceVersion": "2817286", "labels": {"kubernetes.io/metadata.name": "argocd-configmap-81"}, "managedFields": [], "creationTimestamp": "2025-05-31T18:01:05.459475Z"}}',
'2025-05-31T18:01:05.459475Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'3a2fa802-7d63-4d90-9650-066ccb2b34d5',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'grafana-deployment-71',
'grafana',
'4233423',
'{"kind": "Deployment", "metadata": {"uid": "3a2fa802-7d63-4d90-9650-066ccb2b34d5", "name": "grafana-deployment-71", "namespace": "grafana", "resourceVersion": "4233423", "labels": {"kubernetes.io/metadata.name": "grafana-deployment-71"}, "managedFields": [], "creationTimestamp": "2025-06-15T18:01:05.459534Z"}}',
'2025-06-15T18:01:05.459534Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'cb0b2f18-4f8c-4a4c-8fc3-51dd6a78f1ae',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'prod-pod-70',
'prod',
'3004353',
'{"kind": "Pod", "metadata": {"uid": "cb0b2f18-4f8c-4a4c-8fc3-51dd6a78f1ae", "name": "prod-pod-70", "namespace": "prod", "resourceVersion": "3004353", "labels": {"kubernetes.io/metadata.name": "prod-pod-70"}, "managedFields": [], "creationTimestamp": "2025-06-05T18:01:05.459603Z"}}',
'2025-06-05T18:01:05.459603Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'74ce32a5-ed51-4bdf-87d4-5a3f6baed396',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'grafana-configmap-95',
'grafana',
'9185913',
'{"kind": "ConfigMap", "metadata": {"uid": "74ce32a5-ed51-4bdf-87d4-5a3f6baed396", "name": "grafana-configmap-95", "namespace": "grafana", "resourceVersion": "9185913", "labels": {"kubernetes.io/metadata.name": "grafana-configmap-95"}, "managedFields": [], "creationTimestamp": "2025-06-16T18:01:05.459654Z"}}',
'2025-06-16T18:01:05.459654Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'bd3c5e09-b5d8-417b-967f-9cd79373e25e',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'dev-service-24',
'dev',
'7317991',
'{"kind": "Service", "metadata": {"uid": "bd3c5e09-b5d8-417b-967f-9cd79373e25e", "name": "dev-service-24", "namespace": "dev", "resourceVersion": "7317991", "labels": {"kubernetes.io/metadata.name": "dev-service-24"}, "managedFields": [], "creationTimestamp": "2025-06-05T18:01:05.459707Z"}}',
'2025-06-05T18:01:05.459707Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'3d1f63c7-07bc-47ee-b6b2-e4a7e5e4d089',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'monitoring-deployment-22',
'monitoring',
'7428470',
'{"kind": "Deployment", "metadata": {"uid": "3d1f63c7-07bc-47ee-b6b2-e4a7e5e4d089", "name": "monitoring-deployment-22", "namespace": "monitoring", "resourceVersion": "7428470", "labels": {"kubernetes.io/metadata.name": "monitoring-deployment-22"}, "managedFields": [], "creationTimestamp": "2025-05-17T18:01:05.459763Z"}}',
'2025-05-17T18:01:05.459763Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'3b5f7e83-59de-4ea8-948c-4b82250bdd55',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'argocd-pod-33',
'argocd',
'2486735',
'{"kind": "Pod", "metadata": {"uid": "3b5f7e83-59de-4ea8-948c-4b82250bdd55", "name": "argocd-pod-33", "namespace": "argocd", "resourceVersion": "2486735", "labels": {"kubernetes.io/metadata.name": "argocd-pod-33"}, "managedFields": [], "creationTimestamp": "2025-05-18T18:01:05.459821Z"}}',
'2025-05-18T18:01:05.459821Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'ebbf19a5-c235-4993-b950-0d4fbaf81435',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'grafana-deployment-56',
'grafana',
'8731431',
'{"kind": "Deployment", "metadata": {"uid": "ebbf19a5-c235-4993-b950-0d4fbaf81435", "name": "grafana-deployment-56", "namespace": "grafana", "resourceVersion": "8731431", "labels": {"kubernetes.io/metadata.name": "grafana-deployment-56"}, "managedFields": [], "creationTimestamp": "2025-06-21T18:01:05.459890Z"}}',
'2025-06-21T18:01:05.459890Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'3cea5c7f-dd2a-49d8-a0df-7bf7a872728f',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'logging-deployment-73',
'logging',
'9342620',
'{"kind": "Deployment", "metadata": {"uid": "3cea5c7f-dd2a-49d8-a0df-7bf7a872728f", "name": "logging-deployment-73", "namespace": "logging", "resourceVersion": "9342620", "labels": {"kubernetes.io/metadata.name": "logging-deployment-73"}, "managedFields": [], "creationTimestamp": "2025-06-14T18:01:05.459960Z"}}',
'2025-06-14T18:01:05.459960Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'e4dc5799-be75-47ff-9c71-b474542c6cfd',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'argocd-pod-58',
'argocd',
'5959897',
'{"kind": "Pod", "metadata": {"uid": "e4dc5799-be75-47ff-9c71-b474542c6cfd", "name": "argocd-pod-58", "namespace": "argocd", "resourceVersion": "5959897", "labels": {"kubernetes.io/metadata.name": "argocd-pod-58"}, "managedFields": [], "creationTimestamp": "2025-06-18T18:01:05.460008Z"}}',
'2025-06-18T18:01:05.460008Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'cf047b20-d51a-4167-9ea4-c53afb3e2120',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'team-a-pod-75',
'team-a',
'2874810',
'{"kind": "Pod", "metadata": {"uid": "cf047b20-d51a-4167-9ea4-c53afb3e2120", "name": "team-a-pod-75", "namespace": "team-a", "resourceVersion": "2874810", "labels": {"kubernetes.io/metadata.name": "team-a-pod-75"}, "managedFields": [], "creationTimestamp": "2025-06-21T18:01:05.460069Z"}}',
'2025-06-21T18:01:05.460069Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'2e42c0f8-9adb-44f1-b086-4ce8af9cc099',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'logging-service-20',
'logging',
'4012307',
'{"kind": "Service", "metadata": {"uid": "2e42c0f8-9adb-44f1-b086-4ce8af9cc099", "name": "logging-service-20", "namespace": "logging", "resourceVersion": "4012307", "labels": {"kubernetes.io/metadata.name": "logging-service-20"}, "managedFields": [], "creationTimestamp": "2025-06-20T18:01:05.460119Z"}}',
'2025-06-20T18:01:05.460119Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'4e1a64d9-f0c0-46f4-a460-1238b6fb79cf',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'staging-configmap-83',
'staging',
'7304115',
'{"kind": "ConfigMap", "metadata": {"uid": "4e1a64d9-f0c0-46f4-a460-1238b6fb79cf", "name": "staging-configmap-83", "namespace": "staging", "resourceVersion": "7304115", "labels": {"kubernetes.io/metadata.name": "staging-configmap-83"}, "managedFields": [], "creationTimestamp": "2025-06-04T18:01:05.460179Z"}}',
'2025-06-04T18:01:05.460179Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'd8e7b70d-a079-4146-991c-2e8faecdb7dd',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'ingress-nginx-configmap-79',
'ingress-nginx',
'4345788',
'{"kind": "ConfigMap", "metadata": {"uid": "d8e7b70d-a079-4146-991c-2e8faecdb7dd", "name": "ingress-nginx-configmap-79", "namespace": "ingress-nginx", "resourceVersion": "4345788", "labels": {"kubernetes.io/metadata.name": "ingress-nginx-configmap-79"}, "managedFields": [], "creationTimestamp": "2025-06-24T18:01:05.460232Z"}}',
'2025-06-24T18:01:05.460232Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'8fa64244-e477-4344-b604-6e273d3b74cd',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'dev-deployment-6',
'dev',
'9121912',
'{"kind": "Deployment", "metadata": {"uid": "8fa64244-e477-4344-b604-6e273d3b74cd", "name": "dev-deployment-6", "namespace": "dev", "resourceVersion": "9121912", "labels": {"kubernetes.io/metadata.name": "dev-deployment-6"}, "managedFields": [], "creationTimestamp": "2025-05-15T18:01:05.460283Z"}}',
'2025-05-15T18:01:05.460283Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'f2f14e0a-efd5-4bc6-95e1-c3cf3a9ddf45',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'team-b-service-94',
'team-b',
'9077814',
'{"kind": "Service", "metadata": {"uid": "f2f14e0a-efd5-4bc6-95e1-c3cf3a9ddf45", "name": "team-b-service-94", "namespace": "team-b", "resourceVersion": "9077814", "labels": {"kubernetes.io/metadata.name": "team-b-service-94"}, "managedFields": [], "creationTimestamp": "2025-05-28T18:01:05.460411Z"}}',
'2025-05-28T18:01:05.460411Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'9f9aac15-dd4e-420a-91ea-4af4cc038b17',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'ingress-nginx-deployment-91',
'ingress-nginx',
'5235000',
'{"kind": "Deployment", "metadata": {"uid": "9f9aac15-dd4e-420a-91ea-4af4cc038b17", "name": "ingress-nginx-deployment-91", "namespace": "ingress-nginx", "resourceVersion": "5235000", "labels": {"kubernetes.io/metadata.name": "ingress-nginx-deployment-91"}, "managedFields": [], "creationTimestamp": "2025-05-30T18:01:05.460483Z"}}',
'2025-05-30T18:01:05.460483Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'a939d377-9f50-4d3e-931b-56eaf6fe96f2',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'dev-pod-26',
'dev',
'7441714',
'{"kind": "Pod", "metadata": {"uid": "a939d377-9f50-4d3e-931b-56eaf6fe96f2", "name": "dev-pod-26", "namespace": "dev", "resourceVersion": "7441714", "labels": {"kubernetes.io/metadata.name": "dev-pod-26"}, "managedFields": [], "creationTimestamp": "2025-05-15T18:01:05.460544Z"}}',
'2025-05-15T18:01:05.460544Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'a1329c64-b7f2-4be8-a038-129187ce43af',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'grafana-service-14',
'grafana',
'2156030',
'{"kind": "Service", "metadata": {"uid": "a1329c64-b7f2-4be8-a038-129187ce43af", "name": "grafana-service-14", "namespace": "grafana", "resourceVersion": "2156030", "labels": {"kubernetes.io/metadata.name": "grafana-service-14"}, "managedFields": [], "creationTimestamp": "2025-06-01T18:01:05.460595Z"}}',
'2025-06-01T18:01:05.460595Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'6de091f3-9e60-45f6-869e-bbf47887427d',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'logging-pod-99',
'logging',
'9498290',
'{"kind": "Pod", "metadata": {"uid": "6de091f3-9e60-45f6-869e-bbf47887427d", "name": "logging-pod-99", "namespace": "logging", "resourceVersion": "9498290", "labels": {"kubernetes.io/metadata.name": "logging-pod-99"}, "managedFields": [], "creationTimestamp": "2025-07-05T18:01:05.460738Z"}}',
'2025-07-05T18:01:05.460738Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'adccb8e0-f67d-4b87-9cca-99ecb7c1f66f',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'kube-system-deployment-34',
'kube-system',
'3068840',
'{"kind": "Deployment", "metadata": {"uid": "adccb8e0-f67d-4b87-9cca-99ecb7c1f66f", "name": "kube-system-deployment-34", "namespace": "kube-system", "resourceVersion": "3068840", "labels": {"kubernetes.io/metadata.name": "kube-system-deployment-34"}, "managedFields": [], "creationTimestamp": "2025-06-28T18:01:05.460837Z"}}',
'2025-06-28T18:01:05.460837Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'6f7a5810-a13e-475b-bcb0-d2929cddca5f',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'ingress-nginx-configmap-57',
'ingress-nginx',
'7841799',
'{"kind": "ConfigMap", "metadata": {"uid": "6f7a5810-a13e-475b-bcb0-d2929cddca5f", "name": "ingress-nginx-configmap-57", "namespace": "ingress-nginx", "resourceVersion": "7841799", "labels": {"kubernetes.io/metadata.name": "ingress-nginx-configmap-57"}, "managedFields": [], "creationTimestamp": "2025-07-01T18:01:05.460886Z"}}',
'2025-07-01T18:01:05.460886Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'63eddb8a-7d25-4d9a-90b0-5cfe79e79044',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'monitoring-configmap-5',
'monitoring',
'9983523',
'{"kind": "ConfigMap", "metadata": {"uid": "63eddb8a-7d25-4d9a-90b0-5cfe79e79044", "name": "monitoring-configmap-5", "namespace": "monitoring", "resourceVersion": "9983523", "labels": {"kubernetes.io/metadata.name": "monitoring-configmap-5"}, "managedFields": [], "creationTimestamp": "2025-05-30T18:01:05.460927Z"}}',
'2025-05-30T18:01:05.460927Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'0f6c3e86-bc81-4db0-872c-f4f09ce55c80',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'argocd-service-66',
'argocd',
'3402128',
'{"kind": "Service", "metadata": {"uid": "0f6c3e86-bc81-4db0-872c-f4f09ce55c80", "name": "argocd-service-66", "namespace": "argocd", "resourceVersion": "3402128", "labels": {"kubernetes.io/metadata.name": "argocd-service-66"}, "managedFields": [], "creationTimestamp": "2025-06-27T18:01:05.460977Z"}}',
'2025-06-27T18:01:05.460977Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'8b6567c8-371a-43d4-8828-11ede42923bf',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'prod-pod-67',
'prod',
'6012583',
'{"kind": "Pod", "metadata": {"uid": "8b6567c8-371a-43d4-8828-11ede42923bf", "name": "prod-pod-67", "namespace": "prod", "resourceVersion": "6012583", "labels": {"kubernetes.io/metadata.name": "prod-pod-67"}, "managedFields": [], "creationTimestamp": "2025-05-28T18:01:05.461029Z"}}',
'2025-05-28T18:01:05.461029Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'e38ffafb-91b3-4133-90a0-92a709a12a18',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'team-b-service-44',
'team-b',
'4014289',
'{"kind": "Service", "metadata": {"uid": "e38ffafb-91b3-4133-90a0-92a709a12a18", "name": "team-b-service-44", "namespace": "team-b", "resourceVersion": "4014289", "labels": {"kubernetes.io/metadata.name": "team-b-service-44"}, "managedFields": [], "creationTimestamp": "2025-06-04T18:01:05.461104Z"}}',
'2025-06-04T18:01:05.461104Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'7cb2ea53-e22d-457c-8cb5-8fb2cbfe8da6',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'staging-configmap-72',
'staging',
'5039501',
'{"kind": "ConfigMap", "metadata": {"uid": "7cb2ea53-e22d-457c-8cb5-8fb2cbfe8da6", "name": "staging-configmap-72", "namespace": "staging", "resourceVersion": "5039501", "labels": {"kubernetes.io/metadata.name": "staging-configmap-72"}, "managedFields": [], "creationTimestamp": "2025-05-20T18:01:05.461236Z"}}',
'2025-05-20T18:01:05.461236Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'f91d248f-bdc3-4cba-acd0-fac4744d4238',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Pod',
'staging-pod-46',
'staging',
'9231220',
'{"kind": "Pod", "metadata": {"uid": "f91d248f-bdc3-4cba-acd0-fac4744d4238", "name": "staging-pod-46", "namespace": "staging", "resourceVersion": "9231220", "labels": {"kubernetes.io/metadata.name": "staging-pod-46"}, "managedFields": [], "creationTimestamp": "2025-06-22T18:01:05.461320Z"}}',
'2025-06-22T18:01:05.461320Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'ed2d8be4-91bb-488a-be12-59564984e737',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Service',
'prod-service-25',
'prod',
'2104462',
'{"kind": "Service", "metadata": {"uid": "ed2d8be4-91bb-488a-be12-59564984e737", "name": "prod-service-25", "namespace": "prod", "resourceVersion": "2104462", "labels": {"kubernetes.io/metadata.name": "prod-service-25"}, "managedFields": [], "creationTimestamp": "2025-05-21T18:01:05.461398Z"}}',
'2025-05-21T18:01:05.461398Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'28aef61b-36cb-490f-b138-721ef6cdb8e0',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'ConfigMap',
'default-configmap-21',
'default',
'2231078',
'{"kind": "ConfigMap", "metadata": {"uid": "28aef61b-36cb-490f-b138-721ef6cdb8e0", "name": "default-configmap-21", "namespace": "default", "resourceVersion": "2231078", "labels": {"kubernetes.io/metadata.name": "default-configmap-21"}, "managedFields": [], "creationTimestamp": "2025-06-16T18:01:05.461461Z"}}',
'2025-06-16T18:01:05.461461Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'02d0d17d-a543-4284-bccb-4eef05fc8136',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'team-b-deployment-28',
'team-b',
'7410600',
'{"kind": "Deployment", "metadata": {"uid": "02d0d17d-a543-4284-bccb-4eef05fc8136", "name": "team-b-deployment-28", "namespace": "team-b", "resourceVersion": "7410600", "labels": {"kubernetes.io/metadata.name": "team-b-deployment-28"}, "managedFields": [], "creationTimestamp": "2025-06-30T18:01:05.461511Z"}}',
'2025-06-30T18:01:05.461511Z'
);

INSERT INTO sandbox_object (id, sandbox_id, kind, name, namespace, resource_version, raw, created_at) VALUES (
'da8ec931-338c-459d-8a64-73deaf621480',
'3d6868ed-4a22-4c6b-a740-359d5fc2816d',
'Deployment',
'ingress-nginx-deployment-71',
'ingress-nginx',
'4460775',
'{"kind": "Deployment", "metadata": {"uid": "da8ec931-338c-459d-8a64-73deaf621480", "name": "ingress-nginx-deployment-71", "namespace": "ingress-nginx", "resourceVersion": "4460775", "labels": {"kubernetes.io/metadata.name": "ingress-nginx-deployment-71"}, "managedFields": [], "creationTimestamp": "2025-07-01T18:01:05.461567Z"}}',
'2025-07-01T18:01:05.461567Z'
);


END;
