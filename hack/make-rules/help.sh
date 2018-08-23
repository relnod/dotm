#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# Returns the help text that gets extracted from the top comments of the
# Makefile.
function help::get_help() {
    sed '/###/q' "$ROOT/Makefile" | sed -e 's/#//g' | cut -c 2-
}

# Returns the help text for a given rule $1.
# everything between the rule definition and the keyword "PHONY" will be taken
# here.
function help::get_help_for_rule() {
    local rule=$1

    sed -n -e "/# === $rule ===/,/.PHONY/ p" "$ROOT/Makefile" | \
        grep -v "===" | \
        grep -v "PHONY" | \
        sed -e 's/#//g' | \
        cut -c 2-
}

# Returns all rules extracted from the Makefile.
# Rule comment has to be of the following form:
#      # === {rule} ===
function help::get_rules() {
    grep "===" "$ROOT/Makefile"| sed -e 's/[= #]//g'
}

# Checks if the rule exists.
function help::rule_exists() {
    local rules=$1
    local rule=$2
    local exists=false

    for r in $rules; do
        if [ "$r" = "$rule" ]; then
            exists="true"
        fi
    done

    echo "$exists"
}

# Prints the general help message.
function help::print_help() {
    local rules=$1
    local help
    local rulehelp

    help=$(help::get_help)

    echo "$help"
    echo ""
    echo "Rules:"
    for rule in $rules; do
        rulehelp=$(help::get_help_for_rule "$rule" | head -n 2)

        printf " %s" "$rule"
        for ((i=1;i<=$((20 - ${#rule}));i++)); do
            printf ' '
        done
        printf "%s" "$rulehelp"
        echo ""
    done
    exit 0
}

# Prints the help message for the given rule $1.
function help::print_help_for_rule() {
    local rule=$1
    local help

    help=$(help::get_help_for_rule "$rule")

    echo "$help"
}

function main() {
    local rule=$1
    local rules

    # Get rules from Makefile
    rules=$(help::get_rules)

    # Check if rule is "help".
    # If so print general help.
    if [ "$1" = "help" ]; then
        help::print_help "$rules"
    fi

    # Check if given $rule exists.
    if [ "$(winston::help::rule_exists "$rules" "$rule")" = false ]; then
        echo "Rule $rule not found."
        echo "Run \`make help\` to list all rules."
        exit 1
    fi

    # Print help message for specific rule.
    help::print_help_for_rule "$rule"
}

main "$RULE"
