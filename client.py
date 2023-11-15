# Implement client side of chat
import socket 
import select 
import sys 
import threading

out_message = ""
flag = False
nickname = "You"


def clearMessage(message: str) -> str:
	res = ""
	for i in message:
		if i != "\x00":
			res += i
	return res


def writeCache(s):
	with open("client_cache.txt", "a") as f:
		f.write(s + '\n')


def main():

	f = open("client_cache.txt", "r")
	last_messages = f.readlines()
	if len(last_messages) != 0:
		last_messages = last_messages[:-6:-1]
		for line in last_messages[::-1]:
			print(line, end='')
		print()
	f.close()

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
		read_socket, _, _ = select.select([server], [], [], 0.25)

		if len(read_socket) > 0:
			if read_socket[0] == server:
				in_message = server.recv(2048).decode('utf-8')
				in_message = clearMessage(in_message)
				print("\r", in_message, sep="")
				print("\rYou: ", end="")
				writeCache(in_message)
		else:
			global out_message
			global flag
			if flag:
				server.send(bytes(out_message, 'utf-8'))
				print("\rYou: ", end="")
				writeCache(out_message) 
				out_message = ""
				flag = False

	server.close()

def reading():
	global flag
	global out_message
	while True:
		if flag == False:
			#print(f"{nickname}: ", end='')
			#out_message = sys.stdin.readline()
			out_message = input("")
			#print("\r", nickname, out_message)
			out_message = out_message[:]	
			flag = True
			sys.stdin.flush()



f = open("client_cache.txt", "a")
f.close()
p1 = threading.Thread(target=reading)
p2 = threading.Thread(target=main)
p1.start()
p2.start()
p1.join()
p2.join()
# main()
