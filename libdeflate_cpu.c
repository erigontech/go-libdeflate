/* Bridge: compile platform-specific CPU features in its own translation unit. */
#if defined(__x86_64__) || defined(__i386__)
#include "./libdeflate/lib/x86/cpu_features.c"
#elif defined(__aarch64__) || defined(__arm__)
#include "./libdeflate/lib/arm/cpu_features.c"
#endif
