# Implement client side of chat
import socket 
import select 
import sys 

def main():

	f = open("client_cache.txt", "r+") 

	last_messages = f.readlines()
	if len(last_messages) != 0:
		print(last_messages[:-6:-1])

	# server_ip = input("Enter server's ip:\t")
	# server_port = int(input("Enter server's port\t"))
	server_ip = "127.0.0.1"
	server_port = int(8080)

	server = socket.socket(socket.AF_INET, socket.SOCK_STREAM) 
	try:
		server.connect((server_ip, server_port)) 
	except Exception as e:
		print(e)
		exit()

	while True:
		read_socket, _, _ = select.select([server], [], [], 1)

		if len(read_socket) > 0:
			if read_socket[0] == server:
				in_message = server.recv(2048)
				print(str(in_message))
				f.write(str(in_message) + '\n')
		else:
			#out_message = sys.stdin.readline()
			out_message = input("You: ")
			server.send(bytes(out_message, 'utf-8'))
			f.write(out_message + '\n') 
	
	server.close()


main()
