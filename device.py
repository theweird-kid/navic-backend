import requests
import time
import pika
import threading

# Define the device ID and the URL of the update endpoint
device_id = "NAV-001"
url = f"http://localhost:8080/api/devices/{device_id}/location"

# Starting location (somewhere in India)
start_location = {
    "lat": 28.6139,  # Latitude for New Delhi
    "lng": 77.2090   # Longitude for New Delhi
}

# Function to simulate walking by updating the location in small increments
def simulate_walking(location, step_size=0.0001):
    # Update latitude and longitude by a small step size
    location["lat"] += step_size
    location["lng"] += step_size
    return location

# Function to continuously update the location every second
def update_location():
    current_location = start_location
    while True:
        current_location = simulate_walking(current_location)
        response = requests.put(url, json=current_location)

        if response.status_code == 200:
            print(f"Location updated successfully: {current_location}")
        else:
            print(f"Failed to update location: {response.status_code} - {response.text}")

        time.sleep(1)

# Function to subscribe to the device's message queue and print messages
def subscribe_to_queue():
    connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
    channel = connection.channel()

    channel.queue_declare(queue=device_id)

    def callback(ch, method, properties, body):
        print(f"Received message: {body.decode()}")

    channel.basic_consume(queue=device_id, on_message_callback=callback, auto_ack=True)

    print(f"Waiting for messages in queue '{device_id}'. To exit press CTRL+C")
    channel.start_consuming()

# Run the location update and message subscription in separate threads
if __name__ == "__main__":
    location_thread = threading.Thread(target=update_location)
    queue_thread = threading.Thread(target=subscribe_to_queue)

    location_thread.start()
    queue_thread.start()

    location_thread.join()
    queue_thread.join()
