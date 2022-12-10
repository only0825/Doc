package com.moon.exception;

public class NameFormatException extends RuntimeException {

    public NameFormatException() {
    }

    public NameFormatException(String message) {
        super(message);
    }
}
