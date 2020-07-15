class MessageBag<K extends number | string, T> {
    private items = new Map<K, Array<T>>();

    add(key: K, value: T): void {
        let messages: Array<T> | undefined = this.items.get(key);
        if (messages === undefined) {
            messages = [];
        }

        if (!messages.includes(value)) {
            messages.push(value);
        }

        this.items.set(key, messages);
    }

    has(key: K): boolean {
        return this.items.has(key);
    }

    first(key: K): T | undefined {
        if (this.has(key)) {
            return this.get(key)[0] ?? undefined;
        }
        return undefined;
    }

    get(key: K): Array<T> {
        let result = this.items.get(key);
        if (result === undefined) {
            return [];
        }
        return result;
    }

    size(): number {
        return this.items.size;
    }

    isEmpty(): boolean {
        return this.size() === 0;
    }
}

export default MessageBag;
