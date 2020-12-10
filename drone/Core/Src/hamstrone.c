#include "hamstrone.h"

HAL_StatusTypeDef HAMSTRONE_MessageTransmit(UART_HandleTypeDef* huart, HAMSTRONE_Message* msg, uint32_t timeout){
	int payloadLen = sizeof(msg->Payload);
	int bodySize = (5 + msg->PayloadLength) * sizeof(uint8_t);
	uint8_t* body = malloc(bodySize);
	body[0] = HAMSTRONE_MESSAGE_MARKER;
	body[1] = 3 + msg->PayloadLength;
	body[2] = msg->Verb;
	body[3] = msg->Noun;
	for (int i = 0; i < msg->PayloadLength; i++) {
		body[4+i] = msg->Payload[i];
	}
	body[4 + msg->PayloadLength] = 0; //CRC
	HAL_StatusTypeDef res = HAL_UART_Transmit(huart, body, bodySize, timeout);
	free(body);
	return res;
}