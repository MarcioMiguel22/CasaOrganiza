<script setup lang="ts">
  import { useI18n } from "vue-i18n";
  import { toast } from "@/components/ui/sonner";
  import { Button } from "@/components/ui/button";
  import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
  import BaseContainer from "@/components/Base/Container.vue";
  import type { PantryOverview, PantryProduct, PantryProductCreate } from "~/lib/api/classes/pantry";
  import MdiAlertCircleOutline from "~icons/mdi/alert-circle-outline";
  import MdiCartOutline from "~icons/mdi/cart-outline";
  import MdiPackageVariant from "~icons/mdi/package-variant";
  import MdiPlus from "~icons/mdi/plus";
  import MdiRefresh from "~icons/mdi/refresh";

  const { t } = useI18n();
  const api = useUserApi();

  definePageMeta({
    middleware: ["auth"],
  });

  useHead({
    title: `CasaOrganiza | ${t("menu.pantry")}`,
  });

  const overview = ref<PantryOverview | null>(null);
  const loading = ref(true);
  const saving = ref(false);
  const errorMessage = ref("");
  const showAddForm = ref(false);
  const newProduct = ref<PantryProductCreate>({
    name: "",
    quantity: 1,
    minQuantity: 0,
    unit: "un.",
    expiryDate: "",
    purchasePrice: 0,
  });

  const amount = (product: PantryProduct): string => {
    const value = Number.isInteger(product.quantity) ? product.quantity : product.quantity.toFixed(2);
    return `${value}${product.unit ? ` ${product.unit}` : ""}`;
  };

  const expiryLabel = (product: PantryProduct): string => {
    if (!product.expiryDate) return t("pantry.no_expiry");
    return new Intl.DateTimeFormat(undefined, { dateStyle: "medium" }).format(
      new Date(`${product.expiryDate}T12:00:00`)
    );
  };

  async function loadOverview() {
    loading.value = true;
    errorMessage.value = "";
    try {
      const response = await api.pantry.overview();
      if (response.error) {
        errorMessage.value = t("pantry.load_error");
        overview.value = null;
      } else {
        overview.value = response.data;
      }
    } catch {
      errorMessage.value = t("pantry.load_error");
      overview.value = null;
    } finally {
      loading.value = false;
    }
  }

  async function createProduct() {
    if (!newProduct.value.name.trim()) return;
    saving.value = true;
    try {
      const response = await api.pantry.createProduct({
        ...newProduct.value,
        name: newProduct.value.name.trim(),
      });
      if (response.error) {
        toast.error(t("pantry.save_error"));
        return;
      }
      newProduct.value = {
        name: "",
        quantity: 1,
        minQuantity: 0,
        unit: "un.",
        expiryDate: "",
        purchasePrice: 0,
      };
      showAddForm.value = false;
      await loadOverview();
      toast.success(t("pantry.product_added"));
    } finally {
      saving.value = false;
    }
  }

  async function consume(product: PantryProduct) {
    const response = await api.pantry.consume(product.id);
    if (response.error) {
      toast.error(t("pantry.save_error"));
      return;
    }
    await loadOverview();
  }

  async function toggleShopping(product: PantryProduct) {
    const response = await api.pantry.toggleShopping(product.id);
    if (response.error) {
      toast.error(t("pantry.save_error"));
      return;
    }
    await loadOverview();
  }

  onMounted(() => {
    void loadOverview();
  });
</script>

<template>
  <BaseContainer class="flex flex-col gap-4">
    <header class="flex flex-wrap items-center justify-between gap-3">
      <div>
        <h1 class="text-2xl font-semibold">{{ t("pantry.title") }}</h1>
        <p class="text-sm text-foreground/70">{{ t("pantry.subtitle") }}</p>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" :disabled="loading" @click="loadOverview">
          <MdiRefresh :class="{ 'animate-spin': loading }" />
          {{ t("pantry.refresh") }}
        </Button>
        <Button @click="showAddForm = !showAddForm">
          <MdiPlus />
          {{ t("pantry.add_product") }}
        </Button>
      </div>
    </header>

    <Card v-if="showAddForm">
      <CardHeader>
        <CardTitle>{{ t("pantry.add_product") }}</CardTitle>
        <CardDescription>{{ t("pantry.add_product_description") }}</CardDescription>
      </CardHeader>
      <CardContent>
        <form class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-6" @submit.prevent="createProduct">
          <label class="flex flex-col gap-1 lg:col-span-2">
            <span class="text-sm font-medium">{{ t("pantry.product_name") }}</span>
            <input v-model="newProduct.name" required class="h-9 rounded-md border bg-background px-3 text-sm" />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium">{{ t("pantry.quantity") }}</span>
            <input
              v-model.number="newProduct.quantity"
              type="number"
              min="0"
              step="0.01"
              class="h-9 rounded-md border bg-background px-3 text-sm"
            />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium">{{ t("pantry.minimum") }}</span>
            <input
              v-model.number="newProduct.minQuantity"
              type="number"
              min="0"
              step="0.01"
              class="h-9 rounded-md border bg-background px-3 text-sm"
            />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium">{{ t("pantry.unit") }}</span>
            <input v-model="newProduct.unit" class="h-9 rounded-md border bg-background px-3 text-sm" />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium">{{ t("pantry.expiry_date") }}</span>
            <input
              v-model="newProduct.expiryDate"
              type="date"
              class="h-9 rounded-md border bg-background px-3 text-sm"
            />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium">{{ t("pantry.purchase_price") }}</span>
            <input
              v-model.number="newProduct.purchasePrice"
              type="number"
              min="0"
              step="0.01"
              class="h-9 rounded-md border bg-background px-3 text-sm"
            />
          </label>
          <div class="flex items-end lg:col-span-6">
            <Button type="submit" :disabled="saving || !newProduct.name.trim()">{{ t("global.save") }}</Button>
          </div>
        </form>
      </CardContent>
    </Card>

    <p v-if="loading" class="text-sm text-foreground/70">{{ t("global.loading") }}</p>
    <p v-else-if="errorMessage" class="rounded-md border border-destructive/40 p-4 text-sm text-destructive">
      {{ errorMessage }}
    </p>

    <template v-else-if="overview">
      <div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
        <Card>
          <CardContent class="flex items-center gap-3 p-4">
            <MdiPackageVariant class="size-7 text-primary" />
            <div>
              <p class="text-sm text-foreground/70">{{ t("pantry.low_stock") }}</p>
              <p class="text-2xl font-semibold">{{ overview.lowStock.length }}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="flex items-center gap-3 p-4">
            <MdiAlertCircleOutline class="size-7 text-amber-600" />
            <div>
              <p class="text-sm text-foreground/70">{{ t("pantry.expiry_alerts") }}</p>
              <p class="text-2xl font-semibold">{{ overview.expiring.length }}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="flex items-center gap-3 p-4">
            <MdiCartOutline class="size-7 text-primary" />
            <div>
              <p class="text-sm text-foreground/70">{{ t("pantry.shopping_list") }}</p>
              <p class="text-2xl font-semibold">{{ overview.shoppingList.length }}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{{ t("pantry.products") }}</CardTitle>
          <CardDescription>{{ t("pantry.products_description") }}</CardDescription>
        </CardHeader>
        <CardContent>
          <p v-if="overview.products.length === 0" class="text-sm text-foreground/70">{{ t("pantry.empty") }}</p>
          <ul v-else class="divide-y">
            <li
              v-for="product in overview.products"
              :key="product.id"
              class="flex flex-wrap items-center justify-between gap-3 py-3"
            >
              <div class="min-w-48 flex-1">
                <p class="font-medium">{{ product.name }}</p>
                <p class="text-sm text-foreground/70">
                  {{ amount(product) }} · {{ t("pantry.expires") }}: {{ expiryLabel(product) }}
                </p>
              </div>
              <div class="flex flex-wrap gap-2">
                <Button size="sm" variant="outline" @click="consume(product)">− 1 {{ product.unit || "un." }}</Button>
                <Button
                  size="sm"
                  :variant="product.onShoppingList ? 'secondary' : 'outline'"
                  @click="toggleShopping(product)"
                >
                  <MdiCartOutline />
                  {{ product.onShoppingList ? t("pantry.remove_from_shopping") : t("pantry.add_to_shopping") }}
                </Button>
              </div>
            </li>
          </ul>
        </CardContent>
      </Card>
    </template>
  </BaseContainer>
</template>
