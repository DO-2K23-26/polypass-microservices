#!/usr/bin/env python3

import json
import time
import uuid
from datetime import datetime
from confluent_kafka import Producer
import random

# Configuration du producer
producer_config = {
    'bootstrap.servers': '0.0.0.0:19092',
    'client.id': 'python-producer'
}

# Création d'un producer
producer = Producer(producer_config)

# Fonction pour générer un événement de tag
def generate_tag_event():
    return {
        'tagId': str(uuid.uuid4()),
        'name': f'Tag-{uuid.uuid4().hex[:8]}',
        'color': '#FF0000',
        'organizationId': str(uuid.uuid4()),
        'timestamp': int(time.time() * 1000)
    }

# Fonction pour générer un événement de dossier
def generate_folder_event():
    return {
        'folderId': str(uuid.uuid4()),
        'name': f'Folder-{uuid.uuid4().hex[:8]}',
        'organizationId': str(uuid.uuid4()),
        'parentId': str(uuid.uuid4()) if random.random() > 0.5 else None,
        'timestamp': int(time.time() * 1000)
    }

# Fonction pour générer un événement d'utilisateur
def generate_user_event():
    return {
        'userId': str(uuid.uuid4()),
        'email': f'user_{uuid.uuid4().hex[:8]}@example.com',
        'organizationId': str(uuid.uuid4()),
        'timestamp': int(time.time() * 1000)
    }

# Fonction pour générer un événement de credential
def generate_credential_event():
    return {
        'credentialId': str(uuid.uuid4()),
        'name': f'Credential-{uuid.uuid4().hex[:8]}',
        'username': f'user_{uuid.uuid4().hex[:4]}',
        'password': f'pass_{uuid.uuid4().hex[:8]}',
        'url': f'https://example.com/{uuid.uuid4().hex[:8]}',
        'folderId': str(uuid.uuid4()),
        'organizationId': str(uuid.uuid4()),
        'timestamp': int(time.time() * 1000)
    }

def delivery_report(err, msg):
    if err is not None:
        print(f'Message delivery failed: {err}')
    else:
        print(f'Message delivered to {msg.topic()} [{msg.partition()}]')

def main():
    # Liste des topics et leurs fonctions de génération d'événements
    topics = {
        'tag_creation': generate_tag_event,
        'tag_deletion': generate_tag_event,
        'tag_update': generate_tag_event,
        'folder_creation': generate_folder_event,
        'folder_deletion': generate_folder_event,
        'folder_update': generate_folder_event,
        'user_creation': generate_user_event,
        'user_deletion': generate_user_event,
        'user_update': generate_user_event,
        'credential_creation': generate_credential_event,
        'credential_deletion': generate_credential_event,
        'credential_update': generate_credential_event
    }

    try:
        while True:
            for topic, generator in topics.items():
                event = generator()
                print(f"Sending event to {topic}: {event}")
                producer.produce(
                    topic=topic,
                    value=json.dumps(event).encode('utf-8'),
                    callback=delivery_report
                )
                producer.poll(0)
            producer.flush()
            time.sleep(5)  # Attendre 5 secondes entre chaque envoi
    except KeyboardInterrupt:
        print("\nStopping producer...")
    finally:
        producer.flush()

if __name__ == "__main__":
    main() 