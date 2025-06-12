#!/usr/bin/env python3

import time
import uuid
from datetime import datetime
from confluent_kafka import Producer
import random
import avro.schema
from avro.io import DatumWriter, BinaryEncoder
import io

# Configuration du producer
producer_config = {
    'bootstrap.servers': '0.0.0.0:19092',
    'client.id': 'python-producer'
}

existing_credentials = {}  # {id: last_known_event_data}

# Création d'un producer
producer = Producer(producer_config)

# Schémas Avro
tag_schema = avro.schema.parse('''
{
	"type": "record",
	"name": "TagEvent",
	"namespace": "com.example",
	"fields": [
		{ "name": "id", "type": "string", "default": "default-id" },
		{ "name": "name", "type": "string", "default": "default-name" },
		{ "name": "color", "type": "string", "default": "#000000" },
		{ "name": "created_at", "type": "string" },
		{ "name": "updated_at", "type": "string" },
		{ "name": "folder_id", "type": "string" },
		{ "name": "created_by", "type": "string" }
	]
}
''')

folder_schema = avro.schema.parse('''
{
	"type": "record",
	"name": "FolderEvent",
	"namespace": "com.example",
	"fields": [
		{ "name": "id", "type": "string", "default": "default-id" },
		{ "name": "name", "type": "string", "default": "default-name" },
		{ "name": "description", "type": "string", "default": "" },
		{ "name": "icon", "type": "string", "default": "" },
		{ "name": "created_at", "type": "string" },
		{ "name": "updated_at", "type": "string" },
		{ "name": "parent_id", "type": ["null", "string"], "default": null },
		{ "name": "members", "type": {"type": "array", "items": "string"}, "default": [] },
		{ "name": "created_by", "type": "string" }
	]
}
''')

user_schema = avro.schema.parse('''
{
    "type": "record",
    "name": "UserEvent",
    "fields": [
        {"name": "userId", "type": "string"},
        {"name": "email", "type": "string"},
        {"name": "timestamp", "type": "long"}
    ]
}
''')

credential_schema = avro.schema.parse('''
{
    "namespace": "com.polypass.events.credential",
    "type": "record",
    "name": "EventCredential",
    "fields": [
        {"name": "credentialId", "type": "string"},
        {"name": "name", "type": ["null", "string"]},
        {"name": "folderId", "type": ["null", "string"]},
        {"name": "tagIds", "type": ["null", {"type": "array", "items": "string"}]},
        {"name": "timestamp", "type": "long"},
        {"name": "created_at", "type": "string"},
        {"name": "updated_at", "type": "string"}
    ]
}
''')

def serialize_avro(schema, data):
    writer = DatumWriter(schema)
    bytes_writer = io.BytesIO()
    encoder = BinaryEncoder(bytes_writer)
    try:
        writer.write(data, encoder)
        return bytes_writer.getvalue()
    except Exception as e:
        print(f"Error serializing data: {e}")
        print(f"Data: {data}")
        print(f"Schema: {schema}")
        raise e

# Fonction pour générer un événement de tag
def generate_tag_event():
    return {
        'id': str(uuid.uuid4()),
        'name': f'Tag-{uuid.uuid4().hex[:8]}',
        'color': f'Color-{uuid.uuid4().hex[:8]}',
        'created_at': datetime.now().isoformat(),
        'updated_at': datetime.now().isoformat(),
        'folder_id': str(uuid.uuid4()),
        'created_by': str(uuid.uuid4()),
    }

# Fonction pour générer un événement de dossier
def generate_folder_event():
    return {
        'id': str(uuid.uuid4()),
        'name': f'Folder-{uuid.uuid4().hex[:8]}',
        'description': f'Description-{uuid.uuid4().hex[:8]}',
        'icon': f'Icon-{uuid.uuid4().hex[:8]}',
        'created_at': datetime.now().isoformat(),
        'updated_at': datetime.now().isoformat(),
        'parent_id': None,  # Avro va automatiquement gérer None comme null
        'members': [str(uuid.uuid4()) for _ in range(2)],  # Génère 2 membres aléatoires
        'created_by': str(uuid.uuid4())
    }

# Fonction pour générer un événement utilisateur
def generate_user_event():
    return {
        'userId': str(uuid.uuid4()),
        'email': f'user_{uuid.uuid4().hex[:8]}@example.com',
        'timestamp': int(time.time() * 1000)
    }

# Fonction pour générer un événement credential
def generate_credential_event():
    credential_id = str(uuid.uuid4())
    data = {
        'credentialId': credential_id,
        'name': f'Credential-{uuid.uuid4().hex[:8]}',
        'folderId': str(uuid.uuid4()),
        'tagIds': [str(uuid.uuid4()) for _ in range(2)],
        'timestamp': int(time.time() * 1000),
        'created_at': datetime.now().isoformat(),
        'updated_at': datetime.now().isoformat()
    }
    existing_credentials[credential_id] = data
    return data


def generate_credential_update_event():
    if not existing_credentials:
        return None  # Rien à mettre à jour
    credential_id = random.choice(list(existing_credentials.keys()))
    old_data = existing_credentials[credential_id]
    updated_data = {
        **old_data,
        'name': f'Credential-Updated-{uuid.uuid4().hex[:6]}' if random.random() < 0.7 else None,
        'folderId': str(uuid.uuid4()) if random.random() < 0.7 else None,
        'tagIds': [str(uuid.uuid4()) for _ in range(random.randint(0, 3))] if random.random() < 0.7 else None,
        'timestamp': int(time.time() * 1000),
        'updated_at': datetime.now().isoformat()
    }
    existing_credentials[credential_id] = updated_data
    return updated_data

def generate_credential_delete_event():
    if not existing_credentials:
        return None
    credential_id = random.choice(list(existing_credentials.keys()))
    deleted_data = {
        'credentialId': credential_id,
        'name': None,
        'folderId': None,
        'tagIds': None,
        'timestamp': int(time.time() * 1000),
        'created_at': existing_credentials[credential_id]['created_at'],
        'updated_at': datetime.now().isoformat()
    }
    del existing_credentials[credential_id]
    return deleted_data

def delivery_report(err, msg):
    if err is not None:
        print(f'Message delivery failed: {err}')
    else:
        print(f'Message delivered to {msg.topic()} [{msg.partition()}]')

def main():
    # Liste des topics et leurs fonctions de génération d'événements
    topics = {
        'folder-creation': (generate_folder_event, folder_schema),
        'tag-creation': (generate_tag_event, tag_schema),
        'credential-creation': (generate_credential_event, credential_schema),
        # 'credential-update': (generate_credential_update_event, credential_schema),
        # 'credential-deletion': (generate_credential_delete_event, credential_schema),
    }

    try:
        while True:
            topic = random.choice(list(topics.keys()))
            generator, schema = topics[topic]
            event = generator()
            if event is None:
                continue
            print(f"Sending event to {topic}: {event}")
            try:
                avro_data = serialize_avro(schema, event)
                producer.produce(topic=topic, value=avro_data, callback=delivery_report)
                producer.poll(0)
            except Exception as e:
                print(f"Error producing message: {e}")
            producer.flush()
            time.sleep(1)

    except KeyboardInterrupt:
        print("\nStopping producer...")
    finally:
        producer.flush()

if __name__ == "__main__":
    main() 