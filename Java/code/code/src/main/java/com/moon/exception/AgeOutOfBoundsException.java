package com.moon.exception;

public class AgeOutOfBoundsException extends RuntimeException {

    public AgeOutOfBoundsException() {
    }

    public AgeOutOfBoundsException(String message) {
        super(message);
    }
}
