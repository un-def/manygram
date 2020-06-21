# manygram

A profile manager for Telegram Desktop app

## Quickstart

1. Create the config:

    ```sh
    manygram config create
    ```

2. Edit the config if necessary.

3. Create a new profile:

    ```sh
    manygram create profile_name
    ```

4. Run Telegram Desktop with the newly created profile:

    ```sh
    manygram run profile_name
    ```

5. Generate a [desktop entry](https://specifications.freedesktop.org/desktop-entry-spec/desktop-entry-spec-latest.html#introduction) for the profile:

    ```sh
    manygram desktop create profile_name
    ```
