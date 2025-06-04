# List available commands
default:
    @just --list

package profile='default':
    nix build \
        --json \
        --no-link \
        --print-build-logs \
        '.#{{ profile }}'
