.PHONY: all daltu scribe daltu.% scribe.%

# Build both modules using their default targets
all: daltu scribe

# Run the default target (first in sub-makefile) for each module
 daltu:
	$(MAKE) -C daltu

 scribe:
	$(MAKE) -C scribe

# Pass-through pattern rules for convenience, e.g.:
#   make daltu.run   -> make -C daltu run
#   make scribe.test -> make -C scribe test
 daltu.%:
	$(MAKE) -C daltu $*

 scribe.%:
	$(MAKE) -C scribe $*
